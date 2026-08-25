package service

import (
	"fmt"

	"reverseproxy/internal/model"
	"reverseproxy/internal/store"
)

type PlanCompiler struct{ cache *store.PlanCache }

func NewPlanCompiler(cache *store.PlanCache) *PlanCompiler { return &PlanCompiler{cache: cache} }

func (c *PlanCompiler) Compile(routeID string, backends []string) (plan *model.CompiledPlan, err error) {
	plan = &model.CompiledPlan{}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("build route plan: %v", recovered)
			_ = c.cache.Put(plan)
		}
	}()
	model.PopulatePlan(plan, routeID, backends)
	return plan, c.cache.Put(plan)
}

func (c *PlanCompiler) PrimaryBackend(routeID string) (string, bool) {
	plan, ok := c.cache.Get(routeID)
	if !ok {
		return "", false
	}
	return plan.Backends[0], true
}
