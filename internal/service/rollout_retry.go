package service

import (
	"errors"

	"reverseproxy/internal/model"
	"reverseproxy/internal/store"
)

// RolloutEffect 抽象对外部动作的调用（如下发配置到数据面）。
//
// Apply 必须是幂等的：相同 key 的重复调用应识别为“已执行”并返回 nil，
// 而不是再次执行外部动作。key 在同一次灰度发布的各次重试之间保持稳定
// （见 model.RolloutAttempt.IdempotencyKey），由外部系统据此去重，
// 保证即便回包丢失、自动重试，外部动作也只真正执行一次。
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

// Retry 执行一次 route 灰度发布，容忍“动作已生效但回包丢失”导致的自动重试。
//
// 它要保证两件事：
//  1. 外部动作只执行一次：V1、V2 共用同一个幂等键（key 不含 Version），
//     外部系统据此去重——V1 若已生效，V2 的 Apply 命中同一 key 返回 nil 且不重复执行。
//  2. 状态保留最新成功版本：V1 的迟到回调只能确认 V1 自身状态，绝不能把更高版本
//     （V2）的 succeeded 退回 running——这一保护由 store 的版本/终态守卫实现
//     （见 versioned_rollout.go），由 ConfirmLate 落地。
//
// releaseOld 在 V1 的迟到回调最终到达时被关闭。Retry 在 V2 成功后阻塞等待它，
// 以便在同一调用里完成 V1 的收尾（ConfirmLate）；若 V1 直接成功则无需等待。
func (c *RolloutRetryCoordinator) Retry(routeID string, releaseOld <-chan struct{}) error {
	first := model.RolloutAttempt{RouteID: routeID, Version: 1, Status: model.RolloutStatusRunning}
	if err := c.store.Save(first); err != nil {
		return err
	}

	// 发起第一次外部动作。回包可能丢失——此时 Apply 返回 error，
	// 但动作在外部可能已经生效。
	if err := c.effect.Apply(first.IdempotencyKey()); err == nil {
		first.Status = model.RolloutStatusSucceeded
		return c.store.Save(first)
	}

	second := model.RolloutAttempt{RouteID: routeID, Version: 2, Status: model.RolloutStatusRunning}
	if err := c.store.Save(second); err != nil {
		return err
	}
	// V2 与 V1 共用同一幂等键。外部系统据此识别为同一动作并去重：
	// 若 V1 实际已生效，则此处 Apply 返回 nil 且不会再次执行外部动作。
	if err := c.effect.Apply(second.IdempotencyKey()); err != nil {
		return err
	}
	second.Status = model.RolloutStatusSucceeded
	if err := c.store.Save(second); err != nil {
		return err
	}

	// V2 已成功。等待 V1 的迟到回调到达后做 V1 的收尾——
	// 无论回调确认 V1 成功还是误报 running，store 的版本守卫都保证 V2 的 succeeded 不被回退。
	<-releaseOld
	return c.ConfirmLate(routeID, 1, model.RolloutStatusSucceeded)
}

// ConfirmLate 处理某版本迟到回调的状态确认。
//
// 它把指定版本标记为给定状态，但受 store 守卫约束：
//   - 老版本（如 V1）无法覆盖更高版本（V2）的 succeeded；
//   - 终态 succeeded 不会被同/更低版本的 running 覆盖。
//
// 因此即便迟到回调携带的是 running（误报未完成），也无法把已成功的更高版本退回 running。
// 当写入因版本/终态守卫被拒绝时，返回 nil——这表示“迟到回调的状态已被更新版本取代”，
// 属于预期行为而非错误。
func (c *RolloutRetryCoordinator) ConfirmLate(routeID string, version int, status string) error {
	late := model.RolloutAttempt{RouteID: routeID, Version: version, Status: status}
	if err := c.store.Save(late); err != nil {
		if errors.Is(err, store.ErrStaleRollout) {
			// 迟到回调的状态已被更高/终态版本取代，符合预期。
			return nil
		}
		return err
	}
	return nil
}
