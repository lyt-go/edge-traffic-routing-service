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
	// Snapshot 已返回独立拷贝，此处再拷贝一次是为了与 store 完全解耦：
	// 即便将来有人修改 Snapshot 的实现，Begin 的语义仍是"冻结 Begin 时刻的权重"。
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

// LastSnapshot 返回 Begin 时刻冻结的权重快照的拷贝，外部对其的任何读写
// 都不会污染聚合器内部状态，也不会影响后续读取。
func (a *WeightAggregator) LastSnapshot() map[string]*model.RouteWeight {
	a.mu.RLock()
	src := a.last
	a.mu.RUnlock()
	out := make(map[string]*model.RouteWeight, len(src))
	for k, v := range src {
		cp := *v
		out[k] = &cp
	}
	return out
}
