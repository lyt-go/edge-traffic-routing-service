// Strategy 领域模型（负载均衡调度策略）。
package model

import (
	"strings"
	"time"
)

const (
	StrategyStatusActive   = "active"
	StrategyStatusInactive = "inactive"
	StrategyTypeRoundRobin = "round-robin"
	StrategyTypeLeastConn  = "least-conn"
	StrategyTypeIpHash     = "ip-hash"
	StrategyTypeRandom     = "random"
)

type Strategy struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (x *Strategy) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "策略名不能为空")
	}
	x.Type = strings.TrimSpace(x.Type)
	if x.Type == "" {
		return NewValidationError("type", "调度类型不能为空")
	}
	if x.Type != "" {
		switch x.Type {
		case "round-robin", "least-conn", "ip-hash", "random":
		default:
			return NewValidationError("type", "调度类型不合法")
		}
	}
	x.Description = strings.TrimSpace(x.Description)
	if x.Status == "" {
		x.Status = StrategyStatusActive
	}
	if !StrategyValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func StrategyValidStatus(s string) bool {
	switch s {
	case "active", "inactive":
		return true
	default:
		return false
	}
}

var strategyTransitions = map[string]map[string]bool{
	StrategyStatusActive:   {"inactive": true},
	StrategyStatusInactive: {"active": true},
}

func StrategyCanTransition(from, to string) bool {
	if m, ok := strategyTransitions[from]; ok {
		return m[to]
	}
	return false
}

type StrategyFilter struct {
	Type   string
	Status string
}

func (f StrategyFilter) Match(x *Strategy) bool {
	if f.Type != "" && x.Type != f.Type {
		return false
	}
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
