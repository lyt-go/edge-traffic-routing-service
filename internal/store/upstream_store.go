package store

import (
	"reverseproxy/internal/model"
)

func (s *MemoryStore) CreateUpstream(x *model.Upstream) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.upstreams {
		if exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.upstreams[x.ID] = x
	return nil
}

func (s *MemoryStore) GetUpstream(id string) (*model.Upstream, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.upstreams[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetUpstreamByName(v string) (*model.Upstream, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.upstreams {
		if x.Name == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListUpstreams() []*model.Upstream {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Upstream, 0, len(s.upstreams))
	for _, x := range s.upstreams {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateUpstream(x *model.Upstream) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.upstreams[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.upstreams {
		if exist.ID != x.ID && exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.upstreams[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteUpstream(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.upstreams[id]; !ok {
		return ErrNotFound
	}
	delete(s.upstreams, id)
	return nil
}
