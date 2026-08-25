package store

import "sync"

type RouteSnapshotCache struct {
	mu      sync.RWMutex
	payload map[string][]byte
}

func NewRouteSnapshotCache() *RouteSnapshotCache {
	return &RouteSnapshotCache{payload: make(map[string][]byte)}
}

func (c *RouteSnapshotCache) Put(routeID string, payload []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.payload[routeID] = payload
}

func (c *RouteSnapshotCache) Get(routeID string) []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.payload[routeID]
}
