package store

import (
	"reverseproxy/internal/model"
)

func (s *MemoryStore) CreateAccessLog(x *model.AccessLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.accessLogs[x.ID] = x
	return nil
}

func (s *MemoryStore) GetAccessLog(id string) (*model.AccessLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x, ok := s.accessLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return x, nil
}

func (s *MemoryStore) ListAccessLogs() []*model.AccessLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AccessLog, 0, len(s.accessLogs))
	for _, x := range s.accessLogs {
		list = append(list, x)
	}
	return list
}

func (s *MemoryStore) DeleteAccessLog(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.accessLogs[id]; !ok {
		return ErrNotFound
	}
	delete(s.accessLogs, id)
	return nil
}
