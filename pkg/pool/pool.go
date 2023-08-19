package pool

import "sync"

type Pool struct {
	pendingConns    map[int]*Connection
	processingConns map[int]*Connection
	lock            sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{
		pendingConns:    make(map[int]*Connection),
		processingConns: make(map[int]*Connection),
	}
}

func (p *Pool) AddConnection(conn *Connection) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pendingConns[conn.ID()] = conn
}

func (p *Pool) MovePendingToProcessing(conn *Connection) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.pendingConns, conn.ID())
	p.processingConns[conn.ID()] = conn
}

func (p *Pool) RemoveConnection(conn *Connection) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.pendingConns, conn.ID())
	delete(p.processingConns, conn.ID())
}

func (p *Pool) GetConnection(id int) *Connection {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if conn, ok := p.pendingConns[id]; ok {
		return conn
	}
	if conn, ok := p.processingConns[id]; ok {
		return conn
	}
	return nil
}

func (p *Pool) GetPendingConnections() []*Connection {
	p.lock.RLock()
	defer p.lock.RUnlock()

	conns := make([]*Connection, 0, len(p.pendingConns))
	for _, conn := range p.pendingConns {
		conns = append(conns, conn)
	}
	return conns
}

func (p *Pool) GetProcessingConnections() []*Connection {
	p.lock.RLock()
	defer p.lock.RUnlock()

	conns := make([]*Connection, 0, len(p.processingConns))
	for _, conn := range p.processingConns {
		conns = append(conns, conn)
	}
	return conns
}
