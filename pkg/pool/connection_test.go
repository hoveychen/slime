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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hoveychen/slime/pkg/token"
)

func TestConnection_Accept(t *testing.T) {
	// Create a new connection.
	conn := &Connection{
		req: make(chan *http.Request),
	}

	// Create a new request.
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Test that Accept returns the request when it is available.
	go func() {
		conn.req <- req
	}()
	if acceptedReq := conn.Accept(context.Background()); acceptedReq != req {
		t.Errorf("Accept() = %v, want %v", acceptedReq, req)
	}

	// Test that Accept returns nil when the context is canceled.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if acceptedReq := conn.Accept(ctx); acceptedReq != nil {
		t.Errorf("Accept() = %v, want nil", acceptedReq)
	}
}

func TestConnection_NewSubmitter(t *testing.T) {
	// Create a new connection.
	conn := &Connection{}

	// Test that NewSubmitter returns an error when the connection is not processing.
	if _, err := conn.NewSubmitter(); err != ErrNotProcessing {
		t.Errorf("NewSubmitter() error = %v, want %v", err, ErrNotProcessing)
	}

	// Set the connection to processing.
	conn.processing.Store(true)

	// Test that NewSubmitter returns the response writer when the connection is processing.
	if _, err := conn.NewSubmitter(); err != nil {
		t.Errorf("NewSubmitter() error = %v, want nil", err)
	}
}

func TestConnection_Delegate(t *testing.T) {
	// Create a new connection.
	conn := &Connection{
		req: make(chan *http.Request),
	}

	// Create a new request and response writer.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	respWriter := httptest.NewRecorder()

	// Test that Delegate returns an error when the connection is already processing.
	conn.processing.Store(true)
	if delegateErr := conn.Delegate(context.Background(), respWriter, req); delegateErr != ErrAlreadyProcessing {
		t.Errorf("Delegate() error = %v, want %v", delegateErr, ErrAlreadyProcessing)
	}
	conn.processing.Store(false)

	// Test that Delegate returns an error when the connection is canceled.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if delegateErr := conn.Delegate(ctx, respWriter, req); delegateErr != context.Canceled {
		t.Errorf("Delegate() error = %v, want %v", delegateErr, context.Canceled)
	}
	conn.processing.Store(false)

	go func() {
		time.Sleep(10 * time.Millisecond)
		conn.Accept(context.Background())
		time.Sleep(10 * time.Millisecond)
		w, _ := conn.NewSubmitter()
		w.Close()
	}()
	// Test that Delegate sets the connection to processing and returns nil when the request is accepted.
	if delegateErr := conn.Delegate(context.Background(), respWriter, req); delegateErr != nil {
		t.Errorf("Delegate() error = %v, want nil", delegateErr)
	}
	if !conn.processing.Load() {
		t.Errorf("Delegate() conn.processing = false, want true")
	}
	conn.processing.Store(false)

	// Test that Delegate returns an error when the context is canceled while waiting for the response writer to finish.
	wait := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Millisecond)
		conn.Accept(context.Background())
		time.Sleep(100 * time.Millisecond)
		_, err := conn.NewSubmitter()
		if err == nil {
			// Expect an error, since the application has been canceled.
			t.Errorf("Delegate() error = %v, want %v", err, context.DeadlineExceeded)
		}
		close(wait)
	}()
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
	if delegateErr := conn.Delegate(ctx, respWriter, req); delegateErr != context.DeadlineExceeded {
		t.Errorf("Delegate() error = %v, want %v", delegateErr, context.DeadlineExceeded)
	}
	<-wait
	cancel()
}

func TestNewConnection(t *testing.T) {
	// Test that NewConnection returns a non-nil connection.
	token := &token.AgentToken{}
	conn := NewConnection(0, token)
	if conn == nil {
		t.Error("NewConnection() returned nil")
	}

	// Test that the connection's agentToken field is set correctly.
	if conn.agentToken != token {
		t.Error("NewConnection() did not set the agentToken field correctly")
	}

	// Test that the connection's req field is a non-nil channel.
	if conn.req == nil {
		t.Error("NewConnection() did not initialize the req field")
	}

	// Test that the connection's id field is a non-zero value.
	if conn.id == 0 {
		t.Error("NewConnection() did not initialize the id field")
	}

	// Test that the connection's since field is a non-zero value.
	if conn.since.IsZero() {
		t.Error("NewConnection() did not initialize the since field")
	}
}

func TestConnection_ID(t *testing.T) {
	// Test that ID returns the correct value.
	conn := &Connection{id: 123}
	if id := conn.ID(); id != 123 {
		t.Errorf("ID() = %d, want %d", id, 123)
	}
}

func TestConnection_Since(t *testing.T) {
	// Test that Since returns the correct value.
	now := time.Now()
	conn := &Connection{since: now}
	if since := conn.Since(); !since.Equal(now) {
		t.Errorf("Since() = %v, want %v", since, now)
	}
}

func TestConnection_AgentID(t *testing.T) {
	// Test that AgentID returns the correct value.
	token := &token.AgentToken{Id: 123}
	conn := &Connection{agentToken: token}
	if id := conn.TokenID(); id != 123 {
		t.Errorf("AgentID() = %d, want %d", id, 123)
	}
}

func TestConnection_AgentName(t *testing.T) {
	// Test that AgentName returns the correct value.
	token := &token.AgentToken{Name: "test"}
	conn := &Connection{agentToken: token}
	if name := conn.AgentName(); name != "test" {
		t.Errorf("AgentName() = %s, want %s", name, "test")
	}
}

func TestConnection_ScopePaths(t *testing.T) {
	// Test that ScopePaths returns the correct value.
	token := &token.AgentToken{ScopePaths: []string{"test"}}
	conn := &Connection{agentToken: token}
	if paths := conn.ScopePaths(); len(paths) != 1 || paths[0] != "test" {
		t.Errorf("ScopePaths() = %v, want %v", paths, []string{"test"})
	}
}

func TestConnection_IsProcessing(t *testing.T) {
	// Test that IsProcessing returns the correct value.
	conn := &Connection{}
	if processing := conn.IsProcessing(); processing {
		t.Errorf("IsProcessing() = %v, want %v", processing, false)
	}
	conn.processing.Store(true)
	if processing := conn.IsProcessing(); !processing {
		t.Errorf("IsProcessing() = %v, want %v", processing, true)
	}
}

func TestConnection_Close_NoResponseWriter(t *testing.T) {
	// Create a new connection.
	conn := &Connection{}

	// Test that Close does not panic when there is no response writer.
	err := errors.New("test error")
	if closeErr := conn.Close(err); closeErr != nil {
		t.Errorf("Close() error = %v, want nil", closeErr)
	}
	if conn.err.Load() != err {
		t.Errorf("Close() did not store the error")
	}
}
