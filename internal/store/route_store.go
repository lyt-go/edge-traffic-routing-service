package store

import (
	"reverseproxy/internal/model"
)

func (s *MemoryStore) CreateRoute(x *model.Route) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.routes {
		if exist.Path == x.Path {
			return ErrConflict
		}
	}
	s.routes[x.ID] = x
	return nil
}

func (s *MemoryStore) GetRoute(id string) (*model.Route, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.routes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetRouteByPath(v string) (*model.Route, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.routes {
		if x.Path == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRoutes() []*model.Route {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Route, 0, len(s.routes))
	for _, x := range s.routes {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateRoute(x *model.Route) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.routes[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.routes {
		if exist.ID != x.ID && exist.Path == x.Path {
			return ErrConflict
		}
	}
	s.routes[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteRoute(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.routes[id]; !ok {
		return ErrNotFound
	}
	delete(s.routes, id)
	return nil
}
