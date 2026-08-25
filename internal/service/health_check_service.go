package service

import (
	"sort"
	"time"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/idgen"
)

func (s *Service) CreateHealthCheck(input model.HealthCheck) (*model.HealthCheck, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetHealthCheckByUpstreamID(input.UpstreamID); err == nil {
		return nil, model.NewValidationError("upstream_id", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateHealthCheck(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetHealthCheck(id string) (*model.HealthCheck, error) {
	return s.store.GetHealthCheck(id)
}

func (s *Service) ListHealthChecks(filter model.HealthCheckFilter, page, size int) ([]*model.HealthCheck, int, error) {
	all := s.store.ListHealthChecks()
	matched := make([]*model.HealthCheck, 0, len(all))
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
		return []*model.HealthCheck{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateHealthCheck(id string, input model.HealthCheck) (*model.HealthCheck, error) {
	exist, err := s.store.GetHealthCheck(id)
	if err != nil {
		return nil, err
	}
	exist.UpstreamID = input.UpstreamID
	exist.Endpoint = input.Endpoint
	exist.IntervalSec = input.IntervalSec
	exist.TimeoutSec = input.TimeoutSec
	exist.FailThreshold = input.FailThreshold
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateHealthCheck(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteHealthCheck(id string) error {
	return s.store.DeleteHealthCheck(id)
}

func (s *Service) TransitionHealthCheck(id, target string) (*model.HealthCheck, error) {
	if !model.HealthCheckValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetHealthCheck(id)
	if err != nil {
		return nil, err
	}
	if !model.HealthCheckCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateHealthCheck(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
