/*
Copyright © 2023 Harry C <hoveychen@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package hub

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/hoveychen/slime/pkg/hwinfo"
	"github.com/hoveychen/slime/pkg/pool"
	"github.com/hoveychen/slime/pkg/token"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

type TokenManager interface {
	Encrypt(*token.AgentToken) (string, error)
	Decrypt(string) (*token.AgentToken, error)
}

type ConnectionInfo struct {
	AgentID      int
	Since        time.Time
	AgentName    string
	ScopePaths   []string
	Processing   bool
	HardwareInfo *hwinfo.HWInfo
}

// Hub server is responsible for:
// 1. Manage agents' connections
// 2. Forward applications' request to the right agent
type HubServer struct {
	tokenMgr    TokenManager
	connPool    *pool.Pool
	concurrent  chan struct{}
	mutex       sync.RWMutex
	hwInfos     map[int]*hwinfo.HWInfo
	appPassword string
}

type HubServerOption func(hs *HubServer)

func WithConcurrent(n int) HubServerOption {
	return func(hs *HubServer) {
		hs.concurrent = make(chan struct{}, n)
	}
}

func WithAppPassword(password string) HubServerOption {
	return func(hs *HubServer) {
		hs.appPassword = password
	}
}

func NewHubServer(secret string, opts ...HubServerOption) *HubServer {
	hs := &HubServer{
		tokenMgr: token.NewTokenManager([]byte(secret)),
		connPool: pool.NewPool(),
		hwInfos:  make(map[int]*hwinfo.HWInfo),
	}
	for _, opt := range opts {
		opt(hs)
	}
	return hs
}

func (hs *HubServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("slime-agent-token")
	if token != "" && r.Method == "POST" {
		var handler http.Handler
		switch r.URL.Path {
		case PathJoin:
			handler = http.HandlerFunc(hs.handleAgentJoin)
		case PathLeave:
			handler = http.HandlerFunc(hs.handleAgentLeave)
		case PathAccept:
			handler = http.HandlerFunc(hs.handleAgentAccept)
		case PathSubmit:
			handler = http.HandlerFunc(hs.handleAgentSubmit)
		default:
			hs.error(w, logrus.WithField("remote", r.RemoteAddr), nil, "Unsupport path")
			return
		}
		handler = hs.wrapTokenValidator(handler)
		handler.ServeHTTP(w, r)
		return
	}

	// We take it as an application request.
	hs.handleAppRequest(w, r)
}

func (hs *HubServer) handleAppRequest(w http.ResponseWriter, r *http.Request) {
	if hs.appPassword != "" && r.Header.Get("slime-app-password") != hs.appPassword {
		hs.replyStatus(w, logrus.WithField("remote", r.RemoteAddr), http.StatusUnauthorized, "Unauthorized", "Invalid password")
		return
	}

	if hs.concurrent != nil {
		hs.concurrent <- struct{}{}
		defer func() {
			<-hs.concurrent
		}()
	}

	for r.Context().Err() == nil {
		conns := hs.connPool.GetPendingConnections()
		rand.Shuffle(len(conns), func(i, j int) {
			conns[i], conns[j] = conns[j], conns[i]
		})
		for _, conn := range conns {
			if len(conn.ScopePaths()) > 0 && !slices.Contains(conn.ScopePaths(), r.URL.Path) {
				continue
			}

			if err := conn.Delegate(r.Context(), w, r); err != nil {
				logrus.WithError(err).WithField("remote", r.RemoteAddr).Error("Failed to delegate request")

				if r.Context().Err() != nil {
					// Prevent agent from submitting the result.
					hs.connPool.RemoveConnection(conn)
				}

				if errors.Is(err, pool.ErrRetry) {
					continue
				}
			}

			return
		}

		// No connections meet the request.
		if r.Header.Get("slime-block") == "" {
			hs.replyStatus(w, logrus.WithField("remote", r.RemoteAddr), http.StatusServiceUnavailable, "No available agent", "No available agent")
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (hs *HubServer) error(w http.ResponseWriter, log *logrus.Entry, err error, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
	if log != nil {
		if err != nil {
			log = log.WithError(err)
		}
		log.Error(msg)
	}
}

func (hs *HubServer) replyStatus(w http.ResponseWriter, log *logrus.Entry, statusCode int, statusMsg, explainMsg string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(statusMsg))
	if log != nil {
		log.Warn(explainMsg)
	}
}

func (hs *HubServer) wrapTokenValidator(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agentLog := logrus.WithField("client", r.RemoteAddr)
		encryptedToken := r.Header.Get("slime-agent-token")
		tok, err := hs.tokenMgr.Decrypt(encryptedToken)
		if err != nil {
			hs.replyStatus(w, agentLog, http.StatusUnauthorized, "Unauthorized", "Failed to decrypt token")
			return
		}

		if tok.ExpireAt > 0 {
			expireAt := time.Unix(tok.ExpireAt, 0)
			if time.Now().After(expireAt) {
				hs.replyStatus(w, agentLog, http.StatusUnauthorized, "Unauthorized", "Token expired")
				return
			}
		}

		_, err = strconv.Atoi(r.Header.Get("slime-agent-id"))
		if err != nil {
			hs.replyStatus(w, agentLog, http.StatusUnauthorized, "Unauthorized", "Failed to parse node id")
			return
		}

		r = r.WithContext(token.NewContext(r.Context(), tok))

		h.ServeHTTP(w, r)
	})
}

func (hs *HubServer) GetHWInfo(agentID int) *hwinfo.HWInfo {
	hs.mutex.RLock()
	defer hs.mutex.RUnlock()
	return hs.hwInfos[agentID]
}

func (hs *HubServer) setHWInfo(agentID int, hwInfo *hwinfo.HWInfo) {
	hs.mutex.Lock()
	defer hs.mutex.Unlock()
	if hwInfo == nil {
		delete(hs.hwInfos, agentID)
		return
	}
	hs.hwInfos[agentID] = hwInfo
}

func (hs *HubServer) GetConnectionsInfos() []*ConnectionInfo {
	var connectionsInfos []*ConnectionInfo

	conns := hs.connPool.GetPendingConnections()
	conns = append(conns, hs.connPool.GetProcessingConnections()...)

	for _, conn := range conns {
		connectionsInfos = append(connectionsInfos, &ConnectionInfo{
			AgentName:    conn.AgentName(),
			AgentID:      conn.AgentID(),
			Since:        conn.Since(),
			ScopePaths:   conn.ScopePaths(),
			Processing:   conn.IsProcessing(),
			HardwareInfo: hs.GetHWInfo(conn.AgentID()),
		})
	}
	return connectionsInfos
}

func (hs *HubServer) handleAgentJoin(w http.ResponseWriter, r *http.Request) {
	agentID, _ := strconv.Atoi(r.Header.Get("slime-agent-id"))
	token := token.FromContext(r.Context())
	agentLog := logrus.WithFields(logrus.Fields{
		"remote":  r.RemoteAddr,
		"agent":   token.GetName(),
		"agentID": agentID,
	})

	var hwInfo hwinfo.HWInfo
	json.NewDecoder(r.Body).Decode(&hwInfo)
	hs.setHWInfo(agentID, &hwInfo)
	agentLog.Info("Agent has arrived.")
}

func (hs *HubServer) handleAgentLeave(w http.ResponseWriter, r *http.Request) {
	agentID, _ := strconv.Atoi(r.Header.Get("slime-agent-id"))
	token := token.FromContext(r.Context())
	agentLog := logrus.WithFields(logrus.Fields{
		"remote":  r.RemoteAddr,
		"agent":   token.GetName(),
		"agentID": agentID,
	})
	hs.setHWInfo(agentID, nil)
	agentLog.Info("Agent has left.")
}

func (hs *HubServer) handleAgentAccept(w http.ResponseWriter, r *http.Request) {
	agentID, _ := strconv.Atoi(r.Header.Get("slime-agent-id"))
	token := token.FromContext(r.Context())
	agentLog := logrus.WithFields(logrus.Fields{
		"remote":  r.RemoteAddr,
		"agent":   token.GetName(),
		"agentID": agentID,
	})

	agentLog.Info("Agent is listening...")
	// Check if there is any existing connection by the agent.
	// Terminate the existing connection if there is any.
	existings := hs.connPool.GetPendingConnections()
	existings = append(existings, hs.connPool.GetProcessingConnections()...)

	for _, conn := range existings {
		if conn.AgentID() != agentID {
			continue
		}
		if err := conn.SubmitError(r.Context(), pool.ErrAgentAlreadyConnected); err != nil {
			agentLog.WithError(err).Error("Failed to terminate the connection")
		} else {
			agentLog.WithField("connectionID", conn.ID()).Warn("Agent already connected. Terminating the existing connection.")
		}
	}

	conn := pool.NewConnection(agentID, token)
	hs.connPool.AddConnection(conn)
	// Blocking, wait for a new job
	req := conn.Accept(r.Context())
	if r.Context().Err() != nil || req == nil {
		hs.connPool.RemoveConnection(conn)
		hs.error(w, agentLog, r.Context().Err(), "Agent accept canceled")
		return
	}

	// Got the job. Remove the connection from the pending pool.
	hs.connPool.MovePendingToProcessing(conn)
	w.Header().Set("slime-connection-id", strconv.Itoa(conn.ID()))
	w.WriteHeader(http.StatusOK)
	if err := req.Write(w); err != nil {
		if err := conn.SubmitError(r.Context(), errors.Join(err, pool.ErrRetry)); err != nil {
			agentLog.WithError(err).Error("Failed to submit error")
		}
		hs.connPool.RemoveConnection(conn)
		agentLog.WithError(err).Error("Failed to serialize request")
		return
	}

	agentLog.WithFields(logrus.Fields{
		"path":   req.URL.Path,
		"method": req.Method,
	}).Info("Agent accepted.")
}

func (hs *HubServer) handleAgentSubmit(w http.ResponseWriter, r *http.Request) {
	agentID, _ := strconv.Atoi(r.Header.Get("slime-agent-id"))
	token := token.FromContext(r.Context())
	agentLog := logrus.WithFields(logrus.Fields{
		"remote":  r.RemoteAddr,
		"agent":   token.GetName(),
		"agentID": agentID,
	})

	connectionID, err := strconv.Atoi(r.Header.Get("slime-connection-id"))
	if err != nil {
		hs.error(w, agentLog, err, "Invalid connection ID")
		return
	}
	agentLog = agentLog.WithField("connection_id", connectionID)

	conn := hs.connPool.GetConnection(connectionID)
	if conn == nil {
		hs.error(w, agentLog, nil, "Connection not found")
		return
	}
	if conn.TokenID() != token.GetId() {
		hs.error(w, agentLog, nil, "Agent ID mismatch")
		return
	}
	hs.connPool.RemoveConnection(conn)

	errorMsg := r.Header.Get("slime-upstream-result")
	if errorMsg != "" {
		if err := conn.SubmitError(r.Context(), errors.New(errorMsg)); err != nil {
			hs.error(w, agentLog, err, "Failed to submit error result")
			return
		}
		// For the agent, this is not an error, just a result.
		w.WriteHeader(http.StatusOK)
		return
	}

	agentLog.Info("Agent is submitting...")
	submitter, err := conn.NewSubmitter()
	if err != nil {
		hs.error(w, agentLog, err, "Failed to submit")
		return
	}
	defer submitter.Close()

	upResp, err := http.ReadResponse(bufio.NewReader(r.Body), nil)
	if err != nil {
		if err := conn.SubmitError(r.Context(), err); err != nil {
			agentLog.WithError(err).Error("Failed to submit error")
		}
		hs.error(w, agentLog, err, "Read upstream response")
		return
	}

	for k, v := range upResp.Header {
		submitter.Header()[k] = v
	}
	submitter.WriteHeader(upResp.StatusCode)
	io.Copy(submitter, upResp.Body)

	agentLog.WithField("contentLength", upResp.ContentLength).Info("Agent submitted.")
	w.WriteHeader(http.StatusOK)
}
