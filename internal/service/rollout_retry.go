package service

import (
	"reverseproxy/internal/model"
	"reverseproxy/internal/store"
)

type RolloutEffect interface {
	Apply(string) error
}

type RolloutRetryCoordinator struct {
	store  *store.VersionedRolloutStore
	effect RolloutEffect
}

func NewRolloutRetryCoordinator(st *store.VersionedRolloutStore, effect RolloutEffect) *RolloutRetryCoordinator {
	return &RolloutRetryCoordinator{store: st, effect: effect}
}

func (c *RolloutRetryCoordinator) Retry(routeID string, releaseOld <-chan struct{}) error {
	first := model.RolloutAttempt{RouteID: routeID, Version: 1, Status: "running"}
	if err := c.store.Save(first); err != nil {
		return err
	}
	if err := c.effect.Apply(first.IdempotencyKey()); err == nil {
		first.Status = "succeeded"
		return c.store.Save(first)
	}
	second := model.RolloutAttempt{RouteID: routeID, Version: 2, Status: "running"}
	if err := c.store.Save(second); err != nil {
		return err
	}
	if err := c.effect.Apply(second.IdempotencyKey()); err != nil {
		return err
	}
	second.Status = "succeeded"
	if err := c.store.Save(second); err != nil {
		return err
	}
	<-releaseOld
	_ = c.store.Save(first)
	return nil
}
