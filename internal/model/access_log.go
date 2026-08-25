// AccessLog 领域模型（访问日志）。
package model

import (
	"strings"
	"time"
)

type AccessLog struct {
	ID         string    `json:"id"`
	RouteID    string    `json:"route_id"`
	UpstreamID string    `json:"upstream_id"`
	ClientIP   string    `json:"client_ip"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status_code"`
	LatencyMs  int64     `json:"latency_ms"`
	Timestamp  int64     `json:"timestamp"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (x *AccessLog) Validate() error {
	x.RouteID = strings.TrimSpace(x.RouteID)
	if x.RouteID == "" {
		return NewValidationError("route_id", "路由ID不能为空")
	}
	x.UpstreamID = strings.TrimSpace(x.UpstreamID)
	if x.UpstreamID == "" {
		return NewValidationError("upstream_id", "上游ID不能为空")
	}
	x.ClientIP = strings.TrimSpace(x.ClientIP)
	x.Method = strings.TrimSpace(x.Method)
	x.Path = strings.TrimSpace(x.Path)
	if x.StatusCode <= 0 {
		return NewValidationError("status_code", "状态码必须为正数")
	}
	if x.LatencyMs < 0 {
		return NewValidationError("latency_ms", "耗时毫秒不能为负数")
	}
	if x.Timestamp <= 0 {
		return NewValidationError("timestamp", "时间戳必须为正数")
	}
	return nil
}

type AccessLogFilter struct {
	RouteId    string
	StatusCode int
	Method     string
}

func (f AccessLogFilter) Match(x *AccessLog) bool {
	if f.RouteId != "" && x.RouteID != f.RouteId {
		return false
	}
	if f.StatusCode != 0 && x.StatusCode != f.StatusCode {
		return false
	}
	if f.Method != "" && x.Method != f.Method {
		return false
	}
	return true
}
