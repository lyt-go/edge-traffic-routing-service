// HealthCheck 领域模型（健康检查）。
package model

import (
	"strings"
	"time"
)

const (
	HealthCheckStatusEnabled  = "enabled"
	HealthCheckStatusDisabled = "disabled"
)

type HealthCheck struct {
	ID            string    `json:"id"`
	UpstreamID    string    `json:"upstream_id"`
	Endpoint      string    `json:"endpoint"`
	IntervalSec   int       `json:"interval_sec"`
	TimeoutSec    int       `json:"timeout_sec"`
	FailThreshold int       `json:"fail_threshold"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (x *HealthCheck) Validate() error {
	x.UpstreamID = strings.TrimSpace(x.UpstreamID)
	if x.UpstreamID == "" {
		return NewValidationError("upstream_id", "上游ID不能为空")
	}
	x.Endpoint = strings.TrimSpace(x.Endpoint)
	if x.Endpoint == "" {
		return NewValidationError("endpoint", "探活路径不能为空")
	}
	if x.IntervalSec <= 0 {
		return NewValidationError("interval_sec", "间隔秒数必须为正数")
	}
	if x.TimeoutSec <= 0 {
		return NewValidationError("timeout_sec", "超时秒数必须为正数")
	}
	if x.FailThreshold <= 0 {
		return NewValidationError("fail_threshold", "失败阈值必须为正数")
	}
	if x.Status == "" {
		x.Status = HealthCheckStatusEnabled
	}
	if !HealthCheckValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func HealthCheckValidStatus(s string) bool {
	switch s {
	case "enabled", "disabled":
		return true
	default:
		return false
	}
}

var health_checkTransitions = map[string]map[string]bool{
	HealthCheckStatusEnabled:  {"disabled": true},
	HealthCheckStatusDisabled: {"enabled": true},
}

func HealthCheckCanTransition(from, to string) bool {
	if m, ok := health_checkTransitions[from]; ok {
		return m[to]
	}
	return false
}

type HealthCheckFilter struct {
	Status     string
	UpstreamId string
}

func (f HealthCheckFilter) Match(x *HealthCheck) bool {
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	if f.UpstreamId != "" && x.UpstreamID != f.UpstreamId {
		return false
	}
	return true
}
