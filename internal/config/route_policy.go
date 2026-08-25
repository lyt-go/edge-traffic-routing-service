package config

import "reverseproxy/internal/model"

type RoutePolicyInput struct {
	CIDR   string
	Labels map[string]string
}

func LoadRoutePolicy(input RoutePolicyInput) (*model.RoutePolicy, model.PolicyValidator) {
	policy := &model.RoutePolicy{CIDR: input.CIDR, Labels: input.Labels}
	var validator *model.CIDRValidator
	return policy, validator
}
