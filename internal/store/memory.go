package store

import (
	"sync"

	"reverseproxy/internal/model"
)

type MemoryStore struct {
	mu           sync.RWMutex
	upstreams    map[string]*model.Upstream
	routes       map[string]*model.Route
	healthChecks map[string]*model.HealthCheck
	strategies   map[string]*model.Strategy
	accessLogs   map[string]*model.AccessLog
	allowLists   map[string]*model.AllowList
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		upstreams:    make(map[string]*model.Upstream),
		routes:       make(map[string]*model.Route),
		healthChecks: make(map[string]*model.HealthCheck),
		strategies:   make(map[string]*model.Strategy),
		accessLogs:   make(map[string]*model.AccessLog),
		allowLists:   make(map[string]*model.AllowList),
	}
}

var _ Store = (*MemoryStore)(nil)
