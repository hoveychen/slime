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
	conn.processing = true

	// Test that NewSubmitter returns the response writer when the connection is processing.
	if _, err := conn.NewSubmitter(); err != nil {
		t.Errorf("NewSubmitter() error = %v, want nil", err)
	}
}

func TestConnection_SubmitError(t *testing.T) {
	// Create a new connection.
	conn := &Connection{}

	// Test that SubmitError returns an error when the connection is not processing.
	if err := conn.SubmitError(context.Background(), nil); err != ErrNotProcessing {
		t.Errorf("SubmitError() error = %v, want %v", err, ErrNotProcessing)
	}

	// Set the connection to processing.
	conn.processing = true

	// Test that SubmitError sets the error and closes the response writer when the connection is processing.
	respWriter := &WriteCloser{closed: make(chan struct{})}
	conn.respWriter = respWriter
	err := errors.New("test error")
	if submitErr := conn.SubmitError(context.Background(), err); submitErr != nil {
		t.Errorf("SubmitError() error = %v, want nil", submitErr)
	}
	if conn.err != err {
		t.Errorf("SubmitError() conn.err = %v, want %v", conn.err, err)
	}
	if !respWriter.IsClosed() {
		t.Errorf("SubmitError() respWriter.closed = false, want true")
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
	conn.processing = true
	if delegateErr := conn.Delegate(context.Background(), respWriter, req); delegateErr != ErrAlreadyProcessing {
		t.Errorf("Delegate() error = %v, want %v", delegateErr, ErrAlreadyProcessing)
	}
	conn.processing = false

	// Test that Delegate returns an error when the connection is canceled.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if delegateErr := conn.Delegate(ctx, respWriter, req); delegateErr != context.Canceled {
		t.Errorf("Delegate() error = %v, want %v", delegateErr, context.Canceled)
	}
	conn.processing = false

	// TODO(yuheng): Use channel instaed of time-based wait to prevent from the racing condition.
	// Create a channel to signal when the Accept method has been called.
	acceptCalled := make(chan struct{})

	// Create a channel to signal when the NewSubmitter method has been called.
	submitterCalled := make(chan struct{})

	// Create a context with a timeout of 50 milliseconds.
	ctx, cancel = context.WithTimeout(context.Background(), time.Millisecond*50)
	defer cancel()

	// Create a mock response writer and request.
	respWriter = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/", nil)

	// Start a goroutine to call the Accept method and signal when it has been called.
	go func() {
		conn.Accept(ctx)
		close(acceptCalled)
	}()

	// Start a goroutine to call the NewSubmitter method and signal when it has been called.
	go func() {
		w, _ := conn.NewSubmitter()
		w.Close()
		close(submitterCalled)
	}()

	// Call the Delegate method and check that it returns nil and sets the connection to processing.
	if delegateErr := conn.Delegate(ctx, respWriter, req); delegateErr != nil {
		t.Errorf("Delegate() error = %v, want nil", delegateErr)
	}
	if !conn.processing {
		t.Errorf("Delegate() conn.processing = false, want true")
	}

	// Wait for the Accept and NewSubmitter methods to be called.
	<-acceptCalled
	<-submitterCalled

	// Reset the processing flag.
	conn.processing = false

	// Create a channel to signal when the error has been checked.
	errorChecked := make(chan struct{})

	// Start a goroutine to call the Delegate method with a canceled context and check that it returns an error.
	go func() {
		if delegateErr := conn.Delegate(ctx, respWriter, req); delegateErr != context.DeadlineExceeded {
			t.Errorf("Delegate() error = %v, want %v", delegateErr, context.DeadlineExceeded)
		}
		close(errorChecked)
	}()

	// Wait for the error to be checked.
	<-errorChecked
}
