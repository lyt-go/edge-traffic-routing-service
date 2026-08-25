// AllowList 领域模型（IP 白名单/黑名单）。
package model

import (
	"strings"
	"time"
)

const (
	AllowListStatusEnabled  = "enabled"
	AllowListStatusDisabled = "disabled"
	AllowListRuleAllow      = "allow"
	AllowListRuleDeny       = "deny"
)

type AllowList struct {
	ID          string    `json:"id"`
	CIDR        string    `json:"cidr"`
	Rule        string    `json:"rule"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (x *AllowList) Validate() error {
	x.CIDR = strings.TrimSpace(x.CIDR)
	if x.CIDR == "" {
		return NewValidationError("cidr", "网段不能为空")
	}
	x.Rule = strings.TrimSpace(x.Rule)
	if x.Rule == "" {
		return NewValidationError("rule", "规则不能为空")
	}
	if x.Rule != "" {
		switch x.Rule {
		case "allow", "deny":
		default:
			return NewValidationError("rule", "规则不合法")
		}
	}
	x.Description = strings.TrimSpace(x.Description)
	if x.Status == "" {
		x.Status = AllowListStatusEnabled
	}
	if !AllowListValidStatus(x.Status) {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func AllowListValidStatus(s string) bool {
	switch s {
	case "enabled", "disabled":
		return true
	default:
		return false
	}
}

var allow_listTransitions = map[string]map[string]bool{
	AllowListStatusEnabled:  {"disabled": true},
	AllowListStatusDisabled: {"enabled": true},
}

func AllowListCanTransition(from, to string) bool {
	if m, ok := allow_listTransitions[from]; ok {
		return m[to]
	}
	return false
}

type AllowListFilter struct {
	Rule   string
	Status string
}

func (f AllowListFilter) Match(x *AllowList) bool {
	if f.Rule != "" && x.Rule != f.Rule {
		return false
	}
	if f.Status != "" && x.Status != f.Status {
		return false
	}
	return true
}
