package service

import (
	"sync"

	"reverseproxy/internal/model"
	"reverseproxy/internal/store"
)

type WeightAggregator struct {
	store *store.RouteWeightStore
	mu    sync.RWMutex
	last  map[string]*model.RouteWeight
}

func NewWeightAggregator(st *store.RouteWeightStore) *WeightAggregator {
	return &WeightAggregator{store: st}
}

func (a *WeightAggregator) Begin(release <-chan struct{}) <-chan int {
	snapshot := a.store.Snapshot()
	a.mu.Lock()
	a.last = snapshot
	a.mu.Unlock()
	result := make(chan int, 1)
	go func() {
		<-release
		total := 0
		for _, weight := range snapshot {
			total += weight.Weight
		}
		result <- total
		close(result)
	}()
	return result
}

func (a *WeightAggregator) LastSnapshot() map[string]*model.RouteWeight {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.last
}
