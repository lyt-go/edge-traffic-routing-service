package store

import (
	"errors"
	"sync"

	"reverseproxy/internal/model"
)

// ErrStaleRollout 在尝试用更老版本覆盖更新版本的状态时返回。
// 这样可保证迟到的旧回调不会把已经成功的更高版本回退到 running。
var ErrStaleRollout = errors.New("stale rollout update")

type VersionedRolloutStore struct {
	mu     sync.RWMutex
	states map[string]model.RolloutAttempt
}

func NewVersionedRolloutStore() *VersionedRolloutStore {
	return &VersionedRolloutStore{states: make(map[string]model.RolloutAttempt)}
}

// Save 保存一次发布尝试。
//
// 用版本号做乐观并发控制：只有当传入版本 >= 已存版本时才接受更新，
// 这样迟到的旧回调（Version 1）即便声称成功，也无法覆盖更高版本（Version 2）
// 已经落库的 succeeded 状态。
func (s *VersionedRolloutStore) Save(attempt model.RolloutAttempt) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if cur, ok := s.states[attempt.RouteID]; ok {
		// 不允许老版本覆盖新版本，防止迟到回调把状态回退。
		if attempt.Version < cur.Version {
			return ErrStaleRollout
		}
		// 同一版本下，running 不能覆盖已经到达的终态（succeeded）。
		// 终态应是单向流转，避免重复/乱序回调把 succeeded 退回 running。
		if attempt.Version == cur.Version && model.IsRolloutTerminal(cur.Status) && !model.IsRolloutTerminal(attempt.Status) {
			return ErrStaleRollout
		}
	}
	s.states[attempt.RouteID] = attempt
	return nil
}

func (s *VersionedRolloutStore) Get(routeID string) (model.RolloutAttempt, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	attempt, ok := s.states[routeID]
	return attempt, ok
}
