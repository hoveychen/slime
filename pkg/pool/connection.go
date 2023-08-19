package pool

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
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
	processing bool
	err        error
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
	return c.processing
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
	if !c.processing {
		return nil, ErrNotProcessing
	}
	if c.err != nil {
		return nil, c.err
	}
	return c.respWriter, nil
}

func (c *Connection) SubmitError(ctx context.Context, err error) error {
	if !c.processing {
		return ErrNotProcessing
	}
	if c.err != nil {
		return c.err
	}
	c.err = err
	c.respWriter.Close()
	return nil
}

func (c *Connection) Delegate(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	if c.processing {
		return ErrAlreadyProcessing
	}
	c.processing = true

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
		c.err = ctx.Err()
		return ctx.Err()
	case <-c.respWriter.Done():
	}
	if c.err != nil {
		return c.err
	}
	return nil
}
