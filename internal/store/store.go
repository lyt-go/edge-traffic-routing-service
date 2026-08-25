// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"reverseproxy/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	CreateUpstream(x *model.Upstream) error
	GetUpstream(id string) (*model.Upstream, error)
	GetUpstreamByName(v string) (*model.Upstream, error)
	ListUpstreams() []*model.Upstream
	UpdateUpstream(x *model.Upstream) error
	DeleteUpstream(id string) error
	CreateRoute(x *model.Route) error
	GetRoute(id string) (*model.Route, error)
	GetRouteByPath(v string) (*model.Route, error)
	ListRoutes() []*model.Route
	UpdateRoute(x *model.Route) error
	DeleteRoute(id string) error
	CreateHealthCheck(x *model.HealthCheck) error
	GetHealthCheck(id string) (*model.HealthCheck, error)
	GetHealthCheckByUpstreamID(v string) (*model.HealthCheck, error)
	ListHealthChecks() []*model.HealthCheck
	UpdateHealthCheck(x *model.HealthCheck) error
	DeleteHealthCheck(id string) error
	CreateStrategy(x *model.Strategy) error
	GetStrategy(id string) (*model.Strategy, error)
	GetStrategyByName(v string) (*model.Strategy, error)
	ListStrategys() []*model.Strategy
	UpdateStrategy(x *model.Strategy) error
	DeleteStrategy(id string) error
	CreateAccessLog(x *model.AccessLog) error
	GetAccessLog(id string) (*model.AccessLog, error)
	ListAccessLogs() []*model.AccessLog
	DeleteAccessLog(id string) error
	CreateAllowList(x *model.AllowList) error
	GetAllowList(id string) (*model.AllowList, error)
	GetAllowListByCIDR(v string) (*model.AllowList, error)
	ListAllowLists() []*model.AllowList
	UpdateAllowList(x *model.AllowList) error
	DeleteAllowList(id string) error
}
