package store

import (
	"reverseproxy/internal/model"
)

func (s *MemoryStore) CreateAllowList(x *model.AllowList) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.allowLists {
		if exist.CIDR == x.CIDR {
			return ErrConflict
		}
	}
	s.allowLists[x.ID] = x
	return nil
}

func (s *MemoryStore) GetAllowList(id string) (*model.AllowList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.allowLists[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) GetAllowListByCIDR(v string) (*model.AllowList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, x := range s.allowLists {
		if x.CIDR == v {
			return x, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListAllowLists() []*model.AllowList {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AllowList, 0, len(s.allowLists))
	for _, x := range s.allowLists {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) UpdateAllowList(x *model.AllowList) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.allowLists[x.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.allowLists {
		if exist.ID != x.ID && exist.CIDR == x.CIDR {
			return ErrConflict
		}
	}
	s.allowLists[x.ID] = x
	return nil
}

func (s *MemoryStore) DeleteAllowList(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.allowLists[id]; !ok {
		return ErrNotFound
	}
	delete(s.allowLists, id)
	return nil
}
