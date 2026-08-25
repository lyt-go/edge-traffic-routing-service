package service

import (
	"sort"
	"time"

	"reverseproxy/internal/model"
	"reverseproxy/pkg/idgen"
)

func (s *Service) CreateAllowList(input model.AllowList) (*model.AllowList, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetAllowListByCIDR(input.CIDR); err == nil {
		return nil, model.NewValidationError("cidr", "已存在")
	}
	now := time.Now()
	input.ID = idgen.Hex()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := s.store.CreateAllowList(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func (s *Service) GetAllowList(id string) (*model.AllowList, error) {
	return s.store.GetAllowList(id)
}

func (s *Service) ListAllowLists(filter model.AllowListFilter, page, size int) ([]*model.AllowList, int, error) {
	all := s.store.ListAllowLists()
	matched := make([]*model.AllowList, 0, len(all))
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
		return []*model.AllowList{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAllowList(id string, input model.AllowList) (*model.AllowList, error) {
	exist, err := s.store.GetAllowList(id)
	if err != nil {
		return nil, err
	}
	exist.CIDR = input.CIDR
	exist.Rule = input.Rule
	exist.Description = input.Description
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateAllowList(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

func (s *Service) DeleteAllowList(id string) error {
	return s.store.DeleteAllowList(id)
}

func (s *Service) TransitionAllowList(id, target string) (*model.AllowList, error) {
	if !model.AllowListValidStatus(target) {
		return nil, model.NewValidationError("status", "目标状态不合法")
	}
	exist, err := s.store.GetAllowList(id)
	if err != nil {
		return nil, err
	}
	if !model.AllowListCanTransition(exist.Status, target) {
		return nil, model.NewValidationError("status", "不允许从 "+exist.Status+" 流转到 "+target)
	}
	exist.Status = target
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateAllowList(exist); err != nil {
		return nil, err
	}
	return exist, nil
}
