package model

import "fmt"

type CompiledPlan struct {
	RouteID  string
	Backends []string
	Ready    bool
}

func PopulatePlan(plan *CompiledPlan, routeID string, backends []string) {
	plan.RouteID = routeID
	for _, backend := range backends {
		if backend == "" {
			panic("empty backend")
		}
		plan.Backends = append(plan.Backends, backend)
	}
	plan.Ready = true
}

func ValidatePlan(plan *CompiledPlan) error {
	if plan == nil || !plan.Ready || len(plan.Backends) == 0 {
		return fmt.Errorf("route plan is incomplete")
	}
	return nil
}
