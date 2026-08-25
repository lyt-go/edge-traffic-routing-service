package service

import (
	"reverseproxy/internal/model"
	"reverseproxy/internal/store"
)

type RolloutPublisher interface {
	Publish(model.RolloutEvent) error
}

type RolloutAudit interface {
	Record(routeID, result string)
}

type RolloutCommitter struct {
	store     *store.RolloutTxnStore
	publisher RolloutPublisher
	audit     RolloutAudit
}

func NewRolloutCommitter(st *store.RolloutTxnStore, publisher RolloutPublisher, audit RolloutAudit) *RolloutCommitter {
	return &RolloutCommitter{store: st, publisher: publisher, audit: audit}
}

func (c *RolloutCommitter) Apply(routeID, state string) error {
	tx := c.store.Begin(routeID, state)
	c.audit.Record(routeID, "published")
	if err := c.publisher.Publish(model.RolloutEvent{RouteID: routeID, State: state}); err != nil {
		return tx.Finish(err)
	}
	if err := tx.Commit(); err != nil {
		return tx.Finish(err)
	}
	return nil
}
