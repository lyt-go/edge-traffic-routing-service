package store

import (
	"reverseproxy/internal/model"
)

func (s *MemoryStore) CreateStrategy(x *model.Strategy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.strategies {
		if exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.strategies[x.ID] = x
	return nil
}

func (s *MemoryStore) GetStrategy(id string) (*model.Strategy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.strategies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetStrategyByName(v string) (*model.Strategy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.strategies {
		if x.Name == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListStrategys() []*model.Strategy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Strategy, 0, len(s.strategies))
	for _, x := range s.strategies {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateStrategy(x *model.Strategy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.strategies[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.strategies {
		if exist.ID != x.ID && exist.Name == x.Name {
			return ErrConflict
		}
	}
	s.strategies[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteStrategy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.strategies[id]; !ok {
		return ErrNotFound
	}
	delete(s.strategies, id)
	return nil
}
