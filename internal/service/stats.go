package service

import (
	"reverseproxy/internal/model"
)

// Overview 返回全局统计总览：每个实体的总数与按状态分布。
func (s *Service) Overview() map[string]interface{} {
	return map[string]interface{}{
		"upstreams":     s.countByStatus("Upstream", s.store.ListUpstreams()),
		"routes":        s.countByStatus("Route", s.store.ListRoutes()),
		"health-checks": s.countByStatus("HealthCheck", s.store.ListHealthChecks()),
		"strategies":    s.countByStatus("Strategy", s.store.ListStrategys()),
		"access-logs":   s.countByStatus("AccessLog", s.store.ListAccessLogs()),
		"allow-lists":   s.countByStatus("AllowList", s.store.ListAllowLists()),
	}
}

func (s *Service) countByStatus(kind string, list interface{}) map[string]interface{} {
	switch items := list.(type) {
	case []*model.Upstream:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	case []*model.Route:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	case []*model.HealthCheck:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	case []*model.Strategy:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	case []*model.AccessLog:
		return map[string]interface{}{"total": len(items)}
	case []*model.AllowList:
		m := map[string]int{}
		for _, it := range items {
			m[it.Status]++
		}
		return map[string]interface{}{"total": len(items), "by_status": m}
	}
	return map[string]interface{}{"total": 0, "by_status": map[string]int{}}
}
