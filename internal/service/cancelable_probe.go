package service

import (
	"context"

	"reverseproxy/internal/store"
)

type CancelableProbe struct{ scheduler *store.RetryScheduler }

func NewCancelableProbe(scheduler *store.RetryScheduler) *CancelableProbe {
	return &CancelableProbe{scheduler: scheduler}
}

func (p *CancelableProbe) Run(ctx context.Context, upstreamID string, retry <-chan struct{}) <-chan error {
	result := make(chan error)
	p.scheduler.Start(context.Background(), upstreamID, retry, result)
	return result
}

func (p *CancelableProbe) Shutdown(ctx context.Context) error { return p.scheduler.Wait(ctx) }
