package service

import (
	"sort"
	"time"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/idgen"
)

func (s *Service) CreateRoute(input model.Route) (*model.Route, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetRouteByPath(input.Path); err == nil {
		return nil, model.NewValidationError("path", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateRoute(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetRoute(id string) (*model.Route, error) {
	return s.store.GetRoute(id)
}

func (s *Service) ListRoutes(filter model.RouteFilter, page, size int) ([]*model.Route, int, error) {
	all := s.store.ListRoutes()
	matched := make([]*model.Route, 0, len(all))
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
		return []*model.Route{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRoute(id string, input model.Route) (*model.Route, error) {
	exist, err := s.store.GetRoute(id)
	if err != nil {
		return nil, err
	}
	exist.Path = input.Path
	exist.Method = input.Method
	exist.UpstreamIDs = input.UpstreamIDs
	exist.Priority = input.Priority
	exist.StripPrefix = input.StripPrefix
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateRoute(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteRoute(id string) error {
	return s.store.DeleteRoute(id)
}

func (s *Service) TransitionRoute(id, target string) (*model.Route, error) {
	if !model.RouteValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetRoute(id)
	if err != nil {
		return nil, err
	}
	if !model.RouteCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateRoute(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
