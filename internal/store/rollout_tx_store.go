package store

import (
	"errors"
	"sync"
)

var ErrRolloutCommitUnavailable = errors.New("rollout commit unavailable")

type RolloutTxnStore struct {
	mu          sync.Mutex
	states      map[string]string
	failCommits int
}

func NewRolloutTxnStore(failCommits int) *RolloutTxnStore {
	return &RolloutTxnStore{states: make(map[string]string), failCommits: failCommits}
}

func (s *RolloutTxnStore) Begin(routeID, state string) *RolloutTxn {
	return &RolloutTxn{store: s, routeID: routeID, state: state}
}

func (s *RolloutTxnStore) State(routeID string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, ok := s.states[routeID]
	return state, ok
}

type RolloutTxn struct {
	store   *RolloutTxnStore
	routeID string
	state   string
	closed  bool
}

func (t *RolloutTxn) Commit() error {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	if t.closed {
		return errors.New("transaction already closed")
	}
	if t.store.failCommits > 0 {
		t.store.failCommits--
		return ErrRolloutCommitUnavailable
	}
	t.store.states[t.routeID] = t.state
	t.closed = true
	return nil
}

func (t *RolloutTxn) Rollback() error {
	t.closed = true
	return nil
}

func (t *RolloutTxn) Finish(prior error) error {
	if prior != nil {
		return t.Rollback()
	}
	return nil
}
