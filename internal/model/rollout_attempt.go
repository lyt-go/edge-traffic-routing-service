package model

import "fmt"

// RolloutAttempt 描述一次 route 灰度发布尝试。
//
// Version 单调递增，每次重试自增；Status 在 running / succeeded 之间单向流转。
type RolloutAttempt struct {
	RouteID string
	Version int
	Status  string
}

// 发布状态。
const (
	RolloutStatusRunning   = "running"
	RolloutStatusSucceeded = "succeeded"
)

// IsRolloutTerminal 判断是否为终态（成功）。终态不可回退到 running。
func IsRolloutTerminal(status string) bool {
	return status == RolloutStatusSucceeded
}

// IdempotencyKey 返回该 route 发布的幂等键。
//
// 注意：键只绑定 RouteID，不绑定 Version。同一次灰度发布可能因为回包丢失而重试，
// 若每次重试都换一个键，外部系统会把它当成新动作再执行一遍——这会重复下发外部动作。
// 因此键在重试过程中保持稳定，由外部系统依据该键去重，保证“只生效一次”。
func (a RolloutAttempt) IdempotencyKey() string {
	return fmt.Sprintf("rollout:%s", a.RouteID)
}
