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
