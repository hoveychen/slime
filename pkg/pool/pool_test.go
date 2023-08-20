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
	"testing"

	"golang.org/x/exp/slices"
)

func TestPool_AddConnection(t *testing.T) {
	p := NewPool()
	conn := &Connection{id: 1}

	p.AddConnection(conn)

	if len(p.pendingConns) != 1 {
		t.Errorf("AddConnection() failed to add connection to pendingConns")
	}
}

func TestPool_MovePendingToProcessing(t *testing.T) {
	p := NewPool()
	conn := &Connection{id: 1}
	p.AddConnection(conn)

	p.MovePendingToProcessing(conn)

	if len(p.pendingConns) != 0 {
		t.Errorf("MovePendingToProcessing() failed to remove connection from pendingConns")
	}
	if len(p.processingConns) != 1 {
		t.Errorf("MovePendingToProcessing() failed to add connection to processingConns")
	}
}

func TestPool_RemoveConnection(t *testing.T) {
	p := NewPool()
	conn := &Connection{id: 1}
	p.AddConnection(conn)
	p.MovePendingToProcessing(conn)

	p.RemoveConnection(conn)

	if len(p.pendingConns) != 0 {
		t.Errorf("RemoveConnection() failed to remove connection from pendingConns")
	}
	if len(p.processingConns) != 0 {
		t.Errorf("RemoveConnection() failed to remove connection from processingConns")
	}
}

func TestPool_GetConnection(t *testing.T) {
	p := NewPool()
	conn := &Connection{id: 1}
	p.AddConnection(conn)

	if p.GetConnection(1) != conn {
		t.Errorf("GetConnection() failed to get connection from pendingConns")
	}

	p.MovePendingToProcessing(conn)

	if p.GetConnection(1) != conn {
		t.Errorf("GetConnection() failed to get connection from processingConns")
	}

	if p.GetConnection(2) != nil {
		t.Errorf("GetConnection() should return nil for non-existent connection")
	}
}

func TestPool_GetPendingConnections(t *testing.T) {
	p := NewPool()
	conn1 := &Connection{id: 1}
	conn2 := &Connection{id: 2}
	p.AddConnection(conn1)
	p.AddConnection(conn2)

	conns := p.GetPendingConnections()

	if len(conns) != 2 {
		t.Errorf("GetPendingConnections() failed to get all pending connections")
	}
	slices.SortFunc(conns, func(i, j *Connection) int {
		return i.id - j.id
	})
	if conns[0] != conn1 || conns[1] != conn2 {
		t.Errorf("GetPendingConnections() returned incorrect connections")
	}
}

func TestPool_GetProcessingConnections(t *testing.T) {
	p := NewPool()
	conn1 := &Connection{id: 1}
	conn2 := &Connection{id: 2}
	p.AddConnection(conn1)
	p.MovePendingToProcessing(conn1)
	p.MovePendingToProcessing(conn2)

	conns := p.GetProcessingConnections()

	if len(conns) != 2 {
		t.Errorf("GetProcessingConnections() failed to get all processing connections")
	}
	slices.SortFunc(conns, func(i, j *Connection) int {
		return i.id - j.id
	})
	if conns[0] != conn1 || conns[1] != conn2 {
		t.Errorf("GetProcessingConnections() returned incorrect connections")
	}
}
