package model

import (
	"fmt"
	"strings"
)

type RoutePolicy struct {
	CIDR   string
	Labels map[string]string
}

type PolicyValidator interface {
	Validate(*RoutePolicy) error
}

type CIDRValidator struct{}

func (v *CIDRValidator) Validate(policy *RoutePolicy) error {
	if v == nil {
		return nil
	}
	if policy.CIDR != "" && !strings.Contains(policy.CIDR, "/") {
		return fmt.Errorf("invalid CIDR")
	}
	return nil
}
