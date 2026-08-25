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
	result := make(chan error, 1)
	p.scheduler.Start(ctx, upstreamID, retry, result)
	return result
}

// Shutdown aborts any in-flight probes and waits for them to finish or ctx to
// expire. Stopping first guarantees a prompt return even when a probe is
// stuck mid-call.
func (p *CancelableProbe) Shutdown(ctx context.Context) error {
	p.scheduler.Stop()
	return p.scheduler.Wait(ctx)
}
