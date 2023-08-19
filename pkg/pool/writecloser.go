package pool

import (
	"io"
	"net/http"
)

var _ io.WriteCloser = (*WriteCloser)(nil)

type WriteCloser struct {
	http.ResponseWriter
	closed   chan struct{}
	isClosed bool
}

func NewWriteCloser(w http.ResponseWriter) *WriteCloser {
	return &WriteCloser{
		ResponseWriter: w,
		closed:         make(chan struct{}),
	}
}

func (wc *WriteCloser) Close() error {
	if wc.isClosed {
		return nil
	}
	wc.isClosed = true
	close(wc.closed)
	return nil
}

func (wc *WriteCloser) IsClosed() bool {
	return wc.isClosed
}

func (wc *WriteCloser) Done() <-chan struct{} {
	return wc.closed
}
