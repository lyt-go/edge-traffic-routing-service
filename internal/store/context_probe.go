package store

import (
	"context"
	"sync"
)

type ProbeBackend interface {
	Check(context.Context, string) error
}

type ContextProbeStore struct {
	mu      sync.Mutex
	backend ProbeBackend
	ctx     context.Context
}

func NewContextProbeStore(backend ProbeBackend) *ContextProbeStore {
	return &ContextProbeStore{backend: backend}
}

func (s *ContextProbeStore) Check(ctx context.Context, upstreamID string) error {
	s.mu.Lock()
	if s.ctx == nil {
		s.ctx = ctx
	}
	shared := s.ctx
	s.mu.Unlock()
	return s.backend.Check(shared, upstreamID)
}
