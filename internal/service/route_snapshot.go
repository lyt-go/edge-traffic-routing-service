package service

import (
	"reverseproxy/internal/model"
	"reverseproxy/internal/store"
)

type RouteSnapshotService struct{ cache *store.RouteSnapshotCache }

func NewRouteSnapshotService(cache *store.RouteSnapshotCache) *RouteSnapshotService {
	return &RouteSnapshotService{cache: cache}
}

func (s *RouteSnapshotService) Submit(routeID string, frame []byte, release <-chan struct{}) <-chan []byte {
	payload := model.ParseBackendFrame(frame)
	s.cache.Put(routeID, payload)
	exported := make(chan []byte, 1)
	go func() {
		<-release
		exported <- payload
		close(exported)
	}()
	return exported
}

func (s *RouteSnapshotService) Cached(routeID string) []byte { return s.cache.Get(routeID) }
