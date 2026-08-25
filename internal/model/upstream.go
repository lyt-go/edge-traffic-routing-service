// Upstream 领域模型（上游节点）。
package model

import (
	"strings"
	"time"
)

const (
	UpstreamStatusHealthy   = "healthy"
	UpstreamStatusUnhealthy = "unhealthy"
)

type Upstream struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Weight    int       `json:"weight"`
	MaxConns  int       `json:"max_conns"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (x *Upstream) Validate() error {
	x.Name = strings.TrimSpace(x.Name)
	if x.Name == "" {
		return NewValidationError("name", "节点名不能为空")
	}
	x.Host = strings.TrimSpace(x.Host)
	if x.Host == "" {
		return NewValidationError("host", "主机不能为空")
	}
	if x.Port <= 0 {
		return NewValidationError("port", "端口必须为正数")
	}
	if x.Weight <= 0 {
		return NewValidationError("weight", "权重必须为正数")
	}
	if x.MaxConns <= 0 {
		return NewValidationError("max_conns", "最大连接数必须为正数")
	}
	if x.Status == "" {
		x.Status = UpstreamStatusHealthy
	}
	if !UpstreamValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func UpstreamValidStatus(s string) bool {
	switch s {
	case "healthy", "unhealthy":
		return true
	default:
		return false
	}
}

var upstreamTransitions = map[string]map[string]bool{
	UpstreamStatusHealthy:   {"unhealthy": true},
	UpstreamStatusUnhealthy: {"healthy": true},
}

func UpstreamCanTransition(from, to string) bool {
	if m, ok := upstreamTransitions[from]; ok {
		return m[to]
	}
	return false
}

type UpstreamFilter struct {
	Status string
}

func (f UpstreamFilter) Match(x *Upstream) bool {
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
