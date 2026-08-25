package service

import (
	"sort"
	"time"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/idgen"
)

func (s *Service) CreateAccessLog(input model.AccessLog) (*model.AccessLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateAccessLog(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetAccessLog(id string) (*model.AccessLog, error) {
	return s.store.GetAccessLog(id)
}

func (s *Service) ListAccessLogs(filter model.AccessLogFilter, page, size int) ([]*model.AccessLog, int, error) {
	all := s.store.ListAccessLogs()
	matched := make([]*model.AccessLog, 0, len(all))
	for _, x := range all {
		if filter.Match(x) {
			matched = append(matched, x)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AccessLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteAccessLog(id string) error {
	return s.store.DeleteAccessLog(id)
}
