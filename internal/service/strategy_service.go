package service

import (
	"sort"
	"time"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/idgen"
)

func (s *Service) CreateStrategy(input model.Strategy) (*model.Strategy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetStrategyByName(input.Name); err == nil {
		return nil, model.NewValidationError("name", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateStrategy(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetStrategy(id string) (*model.Strategy, error) {
	return s.store.GetStrategy(id)
}

func (s *Service) ListStrategys(filter model.StrategyFilter, page, size int) ([]*model.Strategy, int, error) {
	all := s.store.ListStrategys()
	matched := make([]*model.Strategy, 0, len(all))
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
		return []*model.Strategy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateStrategy(id string, input model.Strategy) (*model.Strategy, error) {
	exist, err := s.store.GetStrategy(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Type = input.Type
	exist.Description = input.Description
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateStrategy(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteStrategy(id string) error {
	return s.store.DeleteStrategy(id)
}

func (s *Service) TransitionStrategy(id, target string) (*model.Strategy, error) {
	if !model.StrategyValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetStrategy(id)
	if err != nil {
		return nil, err
	}
	if !model.StrategyCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateStrategy(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
