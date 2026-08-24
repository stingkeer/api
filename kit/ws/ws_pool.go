package ws

import (
	"sync"
	"sync/atomic"

	"go.aew.app/api.v1/log"
)

type connPool struct {
	mu    sync.RWMutex
	conns map[string]*WSCtx
}

var pool = &connPool{conns: make(map[string]*WSCtx)}

var maxConns atomic.Int32

func init() {
	maxConns.Store(1000)
}

func SetMaxConnections(n int32) {
	maxConns.Store(n)
}

func (p *connPool) set(label string, ws *WSCtx) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if int32(len(p.conns)) >= maxConns.Load() {
		log.Warnf("websocket max connections %d reached", maxConns.Load())
		return false
	}
	if _, exists := p.conns[label]; exists {
		log.Infof("Exist Ctx %s", label)
		return false
	}
	p.conns[label] = ws
	return true
}

func (p *connPool) get(label string) *WSCtx {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.conns[label]
}

func (p *connPool) delete(label string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.conns, label)
}

func setWs(label string, ws *WSCtx) bool {
	return pool.set(label, ws)
}

func GetCtx(label string) *WSCtx {
	return pool.get(label)
}

func CloseAll() {
	pool.mu.RLock()
	copies := make([]*WSCtx, 0, len(pool.conns))
	for _, ws := range pool.conns {
		copies = append(copies, ws)
	}
	pool.mu.RUnlock()
	for _, ws := range copies {
		ws.Close()
	}
}

func ActiveCount() int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	return len(pool.conns)
}
