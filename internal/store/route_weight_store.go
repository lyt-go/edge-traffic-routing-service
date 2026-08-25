package store

import (
	"sync"

	"reverseproxy/internal/model"
)

type RouteWeightStore struct {
	mu      sync.RWMutex
	weights map[string]*model.RouteWeight
}

func NewRouteWeightStore() *RouteWeightStore {
	return &RouteWeightStore{weights: make(map[string]*model.RouteWeight)}
}

func (s *RouteWeightStore) Put(weight model.RouteWeight) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.weights[weight.UpstreamID] = &weight
}

func (s *RouteWeightStore) Snapshot() map[string]*model.RouteWeight {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.weights
}
