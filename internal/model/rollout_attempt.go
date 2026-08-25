package model

import "fmt"

type RolloutAttempt struct {
	RouteID string
	Version int
	Status  string
}

func (a RolloutAttempt) IdempotencyKey() string {
	return fmt.Sprintf("rollout:%s:v%d", a.RouteID, a.Version)
}
