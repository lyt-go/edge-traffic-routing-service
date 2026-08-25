package store

import (
	"errors"
	"sync"

	"reverseproxy/internal/model"
)

var ErrStaleRollout = errors.New("stale rollout update")

type VersionedRolloutStore struct {
	mu     sync.RWMutex
	states map[string]model.RolloutAttempt
}

func NewVersionedRolloutStore() *VersionedRolloutStore {
	return &VersionedRolloutStore{states: make(map[string]model.RolloutAttempt)}
}

func (s *VersionedRolloutStore) Save(attempt model.RolloutAttempt) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[attempt.RouteID] = attempt
	return nil
}

func (s *VersionedRolloutStore) Get(routeID string) (model.RolloutAttempt, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	attempt, ok := s.states[routeID]
	return attempt, ok
}
