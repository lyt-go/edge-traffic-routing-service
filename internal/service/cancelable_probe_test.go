package service_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"reverseproxy/internal/service"
	"reverseproxy/internal/store"
)

type gatedRetryBackend struct {
	mu      sync.Mutex
	calls   int
	started chan struct{}
}

func (b *gatedRetryBackend) Call(context.Context, string) error {
	b.mu.Lock()
	b.calls++
	call := b.calls
	b.mu.Unlock()
	if call == 1 {
		close(b.started)
		return errors.New("temporary probe failure")
	}
	return nil
}

func (b *gatedRetryBackend) count() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.calls
}

func TestShutdownDrainsAfterRequestAbort(t *testing.T) {
	backend := &gatedRetryBackend{started: make(chan struct{})}
	probe := service.NewCancelableProbe(store.NewRetryScheduler(backend))
	ctx, cancel := context.WithCancel(context.Background())
	retry := make(chan struct{})
	_ = probe.Run(ctx, "upstream-slow", retry)
	<-backend.started
	cancel()
	time.Sleep(20 * time.Millisecond)
	close(retry)
	time.Sleep(30 * time.Millisecond)
	if got := backend.count(); got != 1 {
		t.Fatalf("cancelled probe continued retrying: calls=%d", got)
	}
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer shutdownCancel()
	if err := probe.Shutdown(shutdownCtx); err != nil {
		t.Fatalf("shutdown waited on cancelled probe: %v", err)
	}
}
