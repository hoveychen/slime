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
package pool

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/hoveychen/slime/pkg/token"
)

var ErrNotProcessing = errors.New("connection is not processing")
var ErrAlreadyProcessing = errors.New("connection is already processing")
var ErrRetry = errors.New("retry")

type Connection struct {
	agentToken *token.AgentToken
	req        chan (*http.Request)
	since      time.Time
	id         int
	processing atomic.Bool
	err        atomic.Value
	respWriter *WriteCloser
}

func NewConnection(token *token.AgentToken) *Connection {
	return &Connection{
		agentToken: token,
		req:        make(chan (*http.Request)),
		id:         rand.Int(),
		since:      time.Now(),
	}
}

func (c *Connection) ID() int {
	return c.id
}

func (c *Connection) Since() time.Time {
	return c.since
}

func (c *Connection) AgentID() int64 {
	return c.agentToken.GetId()
}

func (c *Connection) AgentName() string {
	return c.agentToken.GetName()
}

func (c *Connection) ScopePaths() []string {
	return c.agentToken.GetScopePaths()
}

func (c *Connection) IsProcessing() bool {
	return c.processing.Load()
}

func (c *Connection) Accept(ctx context.Context) *http.Request {
	select {
	case <-ctx.Done():
		return nil
	case req := <-c.req:
		return req
	}
}

func (c *Connection) NewSubmitter() (*WriteCloser, error) {
	if !c.processing.Load() {
		return nil, ErrNotProcessing
	}
	if c.err.Load() != nil {
		return nil, c.err.Load().(error)
	}
	return c.respWriter, nil
}

func (c *Connection) SubmitError(ctx context.Context, err error) error {
	if !c.processing.Load() {
		return ErrNotProcessing
	}
	if c.err.Load() != nil {
		return c.err.Load().(error)
	}
	c.err.Store(err)
	c.respWriter.Close()
	return nil
}

func (c *Connection) Delegate(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	processing := c.processing.Swap(true)
	if processing {
		return ErrAlreadyProcessing
	}

	c.respWriter = NewWriteCloser(w)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case c.req <- req:
	}
	// The request has been accepted by the agent.
	// Wait for the agent to finish processing the request.
	select {
	case <-ctx.Done():
		// The request has been canceled.
		c.err.Store(ctx.Err())
		return ctx.Err()
	case <-c.respWriter.Done():
	}
	if c.err.Load() != nil {
		return c.err.Load().(error)
	}
	return nil
}
