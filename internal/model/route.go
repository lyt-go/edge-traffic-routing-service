// Route 领域模型（路由规则）。
package model

import (
	"strings"
	"time"
)

const (
	RouteStatusEnabled  = "enabled"
	RouteStatusDisabled = "disabled"
	RouteMethodGET      = "GET"
	RouteMethodPOST     = "POST"
	RouteMethodPUT      = "PUT"
	RouteMethodDELETE   = "DELETE"
	RouteMethodPATCH    = "PATCH"
	RouteMethodANY      = "ANY"
)

type Route struct {
	ID          string    `json:"id"`
	Path        string    `json:"path"`
	Method      string    `json:"method"`
	UpstreamIDs []string  `json:"upstream_ids"`
	Priority    int       `json:"priority"`
	StripPrefix bool      `json:"strip_prefix"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (x *Route) Validate() error {
	x.Path = strings.TrimSpace(x.Path)
	if x.Path == "" {
		return NewValidationError("path", "路径不能为空")
	}
	x.Method = strings.TrimSpace(x.Method)
	if x.Method == "" {
		return NewValidationError("method", "方法不能为空")
	}
	if x.Method != "" {
		switch x.Method {
		case "GET", "POST", "PUT", "DELETE", "PATCH", "ANY":
		default:
			return NewValidationError("method", "方法不合法")
		}
	}
	if x.Priority <= 0 {
		return NewValidationError("priority", "优先级必须为正数")
	}
	if x.Status == "" {
		x.Status = RouteStatusEnabled
	}
	if !RouteValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func RouteValidStatus(s string) bool {
	switch s {
	case "enabled", "disabled":
		return true
	default:
		return false
	}
}

var routeTransitions = map[string]map[string]bool{
	RouteStatusEnabled:  {"disabled": true},
	RouteStatusDisabled: {"enabled": true},
}

func RouteCanTransition(from, to string) bool {
	if m, ok := routeTransitions[from]; ok {
		return m[to]
	}
	return false
}

type RouteFilter struct {
	Method string
	Status string
}

func (f RouteFilter) Match(x *Route) bool {
	if f.Method != "" && x.Method != f.Method {
		return false
	}
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
