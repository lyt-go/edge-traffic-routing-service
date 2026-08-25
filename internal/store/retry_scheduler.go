package store

import (
	"context"
	"sync"
)

type RetryProbeBackend interface {
	Call(context.Context, string) error
}

type RetryScheduler struct {
	backend RetryProbeBackend
	wg      sync.WaitGroup
}

func NewRetryScheduler(backend RetryProbeBackend) *RetryScheduler {
	return &RetryScheduler{backend: backend}
}

func (s *RetryScheduler) Start(ctx context.Context, upstreamID string, retry <-chan struct{}, result chan<- error) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := s.backend.Call(context.Background(), upstreamID)
		if err != nil {
			<-retry
			err = s.backend.Call(context.Background(), upstreamID)
		}
		result <- err
	}()
}

func (s *RetryScheduler) Wait(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
