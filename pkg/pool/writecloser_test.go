package pool

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteCloser_Close(t *testing.T) {
	// Create a new response recorder.
	recorder := httptest.NewRecorder()

	// Create a new WriteCloser.
	wc := NewWriteCloser(recorder)

	// Test that Close sets isClosed to true and closes the closed channel.
	if err := wc.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
	if !wc.isClosed {
		t.Errorf("Close() wc.isClosed = false, want true")
	}
	select {
	case <-wc.closed:
	default:
		t.Errorf("Close() did not close wc.closed channel")
	}

	// Test that Close returns nil when called again.
	if err := wc.Close(); err != nil {
		t.Errorf("Close() error = %v, want nil", err)
	}
}

func TestWriteCloser_IsClosed(t *testing.T) {
	// Create a new response recorder.
	recorder := httptest.NewRecorder()

	// Create a new WriteCloser.
	wc := NewWriteCloser(recorder)

	// Test that IsClosed returns false initially.
	if wc.IsClosed() {
		t.Errorf("IsClosed() = true, want false")
	}

	// Test that IsClosed returns true after Close is called.
	wc.Close()
	if !wc.IsClosed() {
		t.Errorf("IsClosed() = false, want true")
	}
}

func TestWriteCloser_Done(t *testing.T) {
	// Create a new response recorder.
	recorder := httptest.NewRecorder()

	// Create a new WriteCloser.
	wc := NewWriteCloser(recorder)

	// Test that Done returns a channel that is not closed initially.
	select {
	case <-wc.Done():
		t.Errorf("Done() channel is already closed")
	default:
	}

	// Test that Done returns a channel that is closed after Close is called.
	wc.Close()
	select {
	case <-wc.Done():
	default:
		t.Errorf("Done() channel is not closed")
	}
}

func TestWriteCloser_Write(t *testing.T) {
	// Create a new response recorder.
	recorder := httptest.NewRecorder()

	// Create a new WriteCloser.
	wc := NewWriteCloser(recorder)

	// Write some data to the WriteCloser.
	data := []byte("hello, world")
	n, err := wc.Write(data)
	if err != nil {
		t.Errorf("Write() error = %v, want nil", err)
	}
	if n != len(data) {
		t.Errorf("Write() n = %d, want %d", n, len(data))
	}

	// Test that the data was written to the response recorder.
	resp := recorder.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Write() status code = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Write() error reading response body: %v", err)
	}
	if !bytes.Equal(body, data) {
		t.Errorf("Write() response body = %q, want %q", body, data)
	}
}

func TestWriteCloser_WriteHeader(t *testing.T) {
	// Create a new response recorder.
	recorder := httptest.NewRecorder()

	// Create a new WriteCloser.
	wc := NewWriteCloser(recorder)

	// Write a status code to the WriteCloser.
	wc.WriteHeader(http.StatusTeapot)

	// Test that the status code was written to the response recorder.
	resp := recorder.Result()
	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("WriteHeader() status code = %d, want %d", resp.StatusCode, http.StatusTeapot)
	}
}
