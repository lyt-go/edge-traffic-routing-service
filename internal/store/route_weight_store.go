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

// Snapshot 返回当前权重的深拷贝，调用方对其的任何读写都不影响 store
// 内部状态，也不会污染此前已发出的快照。
func (s *RouteWeightStore) Snapshot() map[string]*model.RouteWeight {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]*model.RouteWeight, len(s.weights))
	for k, v := range s.weights {
		cp := *v // 拷贝值，避免外部通过指针篡改 store 内部对象
		out[k] = &cp
	}
	return out
}
