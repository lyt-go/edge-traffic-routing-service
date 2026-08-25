package store

import (
	"sync"

	"reverseproxy/internal/model"
)

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
	return lease
}

func (p *RequestLeasePool) Release(lease *model.RequestLease) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.free = append(p.free, lease)
}
