package store

import (
	"context"
	"sync"
)

type RetryProbeBackend interface {
	Call(context.Context, string) error
}

type RetryScheduler struct {
	backend    RetryProbeBackend
	wg         sync.WaitGroup
	stopCtx    context.Context
	stopCancel context.CancelFunc
}

func NewRetryScheduler(backend RetryProbeBackend) *RetryScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &RetryScheduler{backend: backend, stopCtx: ctx, stopCancel: cancel}
}

// Start launches a probe against upstreamID that observes ctx.
//
// Cancellation is honored at every step: an in-flight Call receives ctx, the
// retry wait races ctx.Done(), and the result send never blocks on a caller
// that has gone away. As a result, cancelling the request aborts the probe and
// suppresses the retry so the upstream is never touched again. The probe is
// also bound to the scheduler's stopCtx, so Shutdown can abort calls that are
// still in-flight when the process is closing down.
func (s *RetryScheduler) Start(ctx context.Context, upstreamID string, retry <-chan struct{}, result chan<- error) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// Derive a call context that is cancelled by either the caller's ctx
		// (request cancelled) or the scheduler's stopCtx (shutdown), whichever
		// fires first.
		callCtx, callCancel := context.WithCancel(ctx)
		stop := context.AfterFunc(s.stopCtx, callCancel)
		defer func() {
			stop()
			callCancel()
		}()

		err := s.backend.Call(callCtx, upstreamID)
		if err != nil {
			// Only retry while the probe is still alive; a cancelled request
			// must not produce a second upstream call.
			select {
			case <-retry:
				err = s.backend.Call(callCtx, upstreamID)
			case <-callCtx.Done():
				err = callCtx.Err()
			}
		}

		// Deliver the result without leaking this goroutine. If the caller has
		// gone away (callCtx done) the receive side may be gone too, so a
		// blocking send could hang forever. Deliver when the channel can
		// accept (buffered, or a waiting receiver); otherwise drop the value.
		select {
		case result <- err:
		case <-callCtx.Done():
			select {
			case result <- err:
			default:
			}
		}
	}()
}

// Stop signals all in-flight probes to abort. Idempotent; safe to call more
// than once. Call it before Wait during shutdown to guarantee a prompt return.
func (s *RetryScheduler) Stop() {
	s.stopCancel()
}

// Wait blocks until every in-flight probe has finished or ctx expires. It does
// not itself cancel work — pair with Stop to abort lingering probes on
// shutdown.
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
