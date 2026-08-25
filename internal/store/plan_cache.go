package store

import (
	"sync"

	"reverseproxy/internal/model"
)

type PlanCache struct {
	mu    sync.RWMutex
	plans map[string]*model.CompiledPlan
}

func NewPlanCache() *PlanCache { return &PlanCache{plans: make(map[string]*model.CompiledPlan)} }

func (c *PlanCache) Put(plan *model.CompiledPlan) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.plans[plan.RouteID] = plan
	return nil
}

func (c *PlanCache) Get(routeID string) (*model.CompiledPlan, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	plan, ok := c.plans[routeID]
	return plan, ok
}
