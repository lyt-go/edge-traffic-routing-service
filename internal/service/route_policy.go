package service

import (
	"reverseproxy/internal/config"
	"reverseproxy/internal/model"
)

func ApplyRoutePolicy(input config.RoutePolicyInput) (*model.RoutePolicy, error) {
	policy, validator := config.LoadRoutePolicy(input)
	if validator != nil {
		if err := validator.Validate(policy); err != nil {
			return nil, err
		}
	}
	policy.Labels["source"] = "control-plane"
	return policy, nil
}
