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
package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hoveychen/slime/pkg/hub"
	"github.com/hoveychen/slime/pkg/hwinfo"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const defaultNumWorker = 1

// Agent server is responsible for:
// 1. Maintain connections to hub
// 2. Forward hub's request to the right upstream
type AgentServer struct {
	numWorker   int
	reportHW    bool
	token       string
	upstreamURL *url.URL
	hubURL      *url.URL
	hwInfo      *hwinfo.HWInfo
	agentID     int
}

type AgentServerOption func(as *AgentServer)

func NewAgentServer(hubAddr, upstreamAddr, token string, opts ...AgentServerOption) (*AgentServer, error) {
	hubURL, err := parseAddr(hubAddr)
	if err != nil {
		return nil, err
	}
	upstreamURL, err := parseAddr(upstreamAddr)
	if err != nil {
		return nil, err
	}

	defaultAgentID, err := getOrCreateAgentID()
	if err != nil {
		return nil, err
	}

	as := &AgentServer{
		token:       token,
		numWorker:   defaultNumWorker,
		hubURL:      hubURL,
		upstreamURL: upstreamURL,
		reportHW:    true,
		agentID:     defaultAgentID,
	}
	for _, opt := range opts {
		opt(as)
	}

	if as.reportHW {
		as.hwInfo = hwinfo.NewHWInfo()
	}

	return as, nil
}

func WithNumWorker(num int) AgentServerOption {
	return func(as *AgentServer) {
		as.numWorker = num
	}
}

func WithReportHardware(report bool) AgentServerOption {
	return func(as *AgentServer) {
		as.reportHW = report
	}
}

func WithAgentID(agentID int) AgentServerOption {
	return func(as *AgentServer) {
		as.agentID = agentID
	}
}

var agentIDFile = "agentID.txt"

// getOrCreateAgentID returns the agent ID. If the agent ID was not generated, it will generate a new random one,
// otherwise, it will return the existing agent ID from a local AID file.
func getOrCreateAgentID() (int, error) {
	// Check if the agentID file exists
	if _, err := os.Stat(agentIDFile); os.IsNotExist(err) {
		// Generate a new agent ID if the file does not exist
		newID := int(rand.Int31())
		// Write the new ID to the file
		err = os.WriteFile(agentIDFile, []byte(fmt.Sprint(newID)), 0644)
		if err != nil {
			return 0, fmt.Errorf("failed to write agent ID to file: %v", err)
		}
		return newID, nil
	}

	// Read the existing ID from the file
	data, err := os.ReadFile(agentIDFile)
	if err != nil {
		return 0, fmt.Errorf("failed to read agent ID from file: %v", err)
	}

	// Convert the read data to an integer
	agentID, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, fmt.Errorf("invalid agent ID in file: %v", err)
	}

	return agentID, nil
}

func parseAddr(addr string) (*url.URL, error) {
	// Either host:port or scheme://[userinfo@]host[:port][/path]
	host, port, err := net.SplitHostPort(addr)
	if err == nil {
		if num, err := strconv.Atoi(port); err == nil && num > 0 && num < 65536 {
			return &url.URL{
				Scheme: "http",
				Host:   net.JoinHostPort(host, port),
			}, nil
		}
	}

	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("invalid scheme")
	}
	if u.Host == "" {
		return nil, errors.New("invalid addr")
	}
	return &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		User:   u.User,
		Path:   u.Path,
	}, nil
}

func (as *AgentServer) Run(ctx context.Context) error {

	grp, ctx := errgroup.WithContext(ctx)
	for i := 0; i < as.numWorker; i++ {
		workerNum := i
		grp.Go(func() error {
			agentID := as.agentID + workerNum
			// connect to hub to ensure the hub address and token is correct.
			if err := as.joinHub(ctx, agentID); err != nil {
				return err
			}
			return as.runWorker(ctx, agentID, workerNum)
		})
	}

	return grp.Wait()
}

func (as *AgentServer) newHubAPIRequest(ctx context.Context, agentID int, apiPath string, reader io.Reader) *http.Request {
	u := *as.hubURL
	u.Path = path.Join(u.Path, apiPath)
	req, _ := http.NewRequestWithContext(ctx, "POST", u.String(), reader)
	req.Header.Set("slime-agent-token", as.token)
	req.Header.Set("slime-agent-id", strconv.Itoa(agentID))
	return req
}

func (as *AgentServer) fixUpstreamRequest(r *http.Request) {
	r.URL.Host = as.upstreamURL.Host
	r.URL.Scheme = as.upstreamURL.Scheme
	r.URL.User = as.upstreamURL.User
	r.Host = ""
	r.RequestURI = ""
}

func (as *AgentServer) joinHub(ctx context.Context, agentID int) error {
	var body io.Reader
	if as.hwInfo != nil {
		json, err := json.Marshal(as.hwInfo)
		if err != nil {
			return err
		}
		body = bufio.NewReader(bytes.NewReader(json))
	}
	req := as.newHubAPIRequest(ctx, agentID, hub.PathJoin, body)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	logrus.Infof("Joined hub: %s", as.hubURL)
	return nil
}

func (as *AgentServer) runWorker(ctx context.Context, agentID int, workerNum int) error {
	log := logrus.WithField("worker", workerNum)
	backoffDuration := time.Second
	for ctx.Err() == nil {
		var connectionID string
		func() error {
			acceptReq := as.newHubAPIRequest(ctx, agentID, hub.PathAccept, nil)
			acceptResp, err := http.DefaultClient.Do(acceptReq)
			if err != nil && (errors.Is(err, io.ErrUnexpectedEOF) || strings.Contains(err.Error(), "unexpected EOF")) {
				// The connection has been accepted by hub and got terminated waiting for a task.
				// Retry immediately.
				log.WithError(err).Warnf("Lost connection... Retry right await")
				time.Sleep(time.Millisecond * 100)
				backoffDuration = time.Second
				return nil
			}
			if err != nil {
				log.WithError(err).Warnf("Listening... Retry in %s", backoffDuration)
				time.Sleep(backoffDuration)
				backoffDuration *= 2
				return nil
			}
			defer acceptResp.Body.Close()
			backoffDuration = time.Second

			if acceptResp.StatusCode != http.StatusOK {
				log.WithField("status_code", acceptResp.StatusCode).Error(acceptResp.Status)
				return nil
			}

			connectionID = acceptResp.Header.Get("slime-connection-id")
			if connectionID == "" {
				log.Error("No connection id in accepted request")
				return nil
			}

			upReq, err := http.ReadRequest(bufio.NewReader(acceptResp.Body))
			if err != nil {
				log.WithError(err).Error("Parsing upstream request")
				return err
			}

			log = log.WithField("path", upReq.URL.Path)

			pr, pw := io.Pipe()

			grp, ctx := errgroup.WithContext(ctx)
			recvHeader := make(chan struct{})
			grp.Go(func() error {
				defer pr.Close()

				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-recvHeader:
				}
				submitReq := as.newHubAPIRequest(ctx, agentID, hub.PathSubmit, pr)
				submitReq.Header.Set("slime-connection-id", connectionID)
				submitResp, err := http.DefaultClient.Do(submitReq)
				if err != nil {
					log.WithError(err).Error("Submit result")
					return err
				}
				defer submitResp.Body.Close()
				if submitResp.StatusCode != http.StatusOK {
					log.WithField("status_code", submitResp.StatusCode).Errorf("Submit result: %s", submitResp.Status)
					return errors.New(submitResp.Status)
				}
				log.Info("Result submitted")
				return nil
			})

			grp.Go(func() error {
				defer pw.Close()
				log.Info("Invoke upstream...")

				as.fixUpstreamRequest(upReq)
				upReq = upReq.WithContext(ctx)
				upResp, err := http.DefaultClient.Do(upReq)
				if err != nil {
					log.WithError(err).Error("Invoke upstream")
					return err
				}
				defer upResp.Body.Close()

				log.WithFields(logrus.Fields{
					"status_code":    upResp.StatusCode,
					"content_length": upResp.ContentLength,
				}).Info("Upstream responsed")

				select {
				case recvHeader <- struct{}{}:
				case <-ctx.Done():
					return ctx.Err()
				}

				if err := upResp.Write(pw); err != nil {
					log.WithError(err).Error("Write upstream response")
					return err
				}
				return nil
			})

			if err := grp.Wait(); err != nil {
				log.WithError(err).Error("Worker error")
				return err
			}
			return nil
		}()
	}

	return ctx.Err()
}
