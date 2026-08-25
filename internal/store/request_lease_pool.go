package store

import (
	"sync"

	"reverseproxy/internal/model"
)

// RequestLeasePool 是 RequestLease 的对象池。
// 借出的对象会跨请求复用，因此 Acquire 时会清空自身状态，
// 调用方仍需保证：异步读取必须在归还前取好快照，不能持有
// 借出的 *RequestLease 在归还后继续访问。
type RequestLeasePool struct {
	mu   sync.Mutex
	free []*model.RequestLease
}

func NewRequestLeasePool() *RequestLeasePool { return &RequestLeasePool{} }

func (p *RequestLeasePool) Acquire() *model.RequestLease {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.free) == 0 {
		return &model.RequestLease{Headers: make(map[string]string)}
	}
	last := len(p.free) - 1
	lease := p.free[last]
	p.free = p.free[:last]
	lease.Reset()
	return lease
}

func (p *RequestLeasePool) Release(lease *model.RequestLease) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.free = append(p.free, lease)
}
