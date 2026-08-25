package service_test

import (
	"errors"
	"sync"
	"testing"

	"reverseproxy/internal/service"
	"reverseproxy/internal/store"
)

type idempotentEffect struct {
	mu   sync.Mutex
	seen map[string]bool
	runs int
}

func (e *idempotentEffect) Apply(key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.seen[key] {
		return nil
	}
	e.seen[key] = true
	e.runs++
	if e.runs == 1 {
		return errors.New("reply lost after effect")
	}
	return nil
}

func TestRetryKeepsSingleEffectAndNewestTerminalState(t *testing.T) {
	st := store.NewVersionedRolloutStore()
	effect := &idempotentEffect{seen: make(map[string]bool)}
	coordinator := service.NewRolloutRetryCoordinator(st, effect)
	releaseOld := make(chan struct{})
	close(releaseOld)
	if err := coordinator.Retry("route-canary", releaseOld); err != nil {
		t.Fatalf("rollout retry failed: %v", err)
	}
	if effect.runs != 1 {
		t.Fatalf("external rollout effect ran more than once: runs=%d", effect.runs)
	}
	state, ok := st.Get("route-canary")
	if !ok || state.Version != 2 || state.Status != "succeeded" {
		t.Fatalf("late first attempt replaced the successful retry: %#v ok=%v", state, ok)
	}
}
