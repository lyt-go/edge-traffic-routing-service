package service

import (
	"sort"
	"time"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/idgen"
)

func (s *Service) CreateUpstream(input model.Upstream) (*model.Upstream, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetUpstreamByName(input.Name); err == nil {
		return nil, model.NewValidationError("name", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateUpstream(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetUpstream(id string) (*model.Upstream, error) {
	return s.store.GetUpstream(id)
}

func (s *Service) ListUpstreams(filter model.UpstreamFilter, page, size int) ([]*model.Upstream, int, error) {
	all := s.store.ListUpstreams()
	matched := make([]*model.Upstream, 0, len(all))
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
		return []*model.Upstream{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateUpstream(id string, input model.Upstream) (*model.Upstream, error) {
	exist, err := s.store.GetUpstream(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Host = input.Host
	exist.Port = input.Port
	exist.Weight = input.Weight
	exist.MaxConns = input.MaxConns
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateUpstream(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteUpstream(id string) error {
	return s.store.DeleteUpstream(id)
}

func (s *Service) TransitionUpstream(id, target string) (*model.Upstream, error) {
	if !model.UpstreamValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetUpstream(id)
	if err != nil {
		return nil, err
	}
	if !model.UpstreamCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateUpstream(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
