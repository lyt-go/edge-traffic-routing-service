package store

import (
	"reverseproxy/internal/model"
)

func (s *MemoryStore) CreateHealthCheck(x *model.HealthCheck) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.healthChecks {
		if exist.UpstreamID == x.UpstreamID {
			return ErrConflict
		}
	}
	s.healthChecks[x.ID] = x
	return nil
}

func (s *MemoryStore) GetHealthCheck(id string) (*model.HealthCheck, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.healthChecks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetHealthCheckByUpstreamID(v string) (*model.HealthCheck, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.healthChecks {
		if x.UpstreamID == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListHealthChecks() []*model.HealthCheck {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.HealthCheck, 0, len(s.healthChecks))
	for _, x := range s.healthChecks {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateHealthCheck(x *model.HealthCheck) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.healthChecks[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.healthChecks {
		if exist.ID != x.ID && exist.UpstreamID == x.UpstreamID {
			return ErrConflict
		}
	}
	s.healthChecks[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteHealthCheck(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.healthChecks[id]; !ok {
		return ErrNotFound
	}
	delete(s.healthChecks, id)
	return nil
}
