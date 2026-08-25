# 反向代理 / 负载均衡

纯 Go 标准库（`net/http`）实现的后端服务，零第三方依赖，标准分层（cmd/internal/pkg），开箱即跑。

## 运行

```bash
cd origin
go run ./cmd/server
# 默认监听 :8080，可用 PORT/ADDR/MAX_PAGE_SIZE 环境变量覆盖
```

## 业务实体

Upstream（上游节点）、Route（路由规则）、HealthCheck（健康检查）、Strategy（负载均衡调度策略）、AccessLog（访问日志）、AllowList（IP 白名单/黑名单）

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`。

| 模块 | 接口 | 说明 |
|------|------|------|
| upstreams | POST /api/upstreams | 创建Upstream |
| upstreams | GET /api/upstreams | 列表查询（分页+筛选） |
| upstreams | GET /api/upstreams/{id} | 详情 |
| upstreams | PUT /api/upstreams/{id} | 更新 |
| upstreams | DELETE /api/upstreams/{id} | 删除 |
| upstreams | PATCH /api/upstreams/{id}/status | 状态流转 |
| routes | POST /api/routes | 创建Route |
| routes | GET /api/routes | 列表查询（分页+筛选） |
| routes | GET /api/routes/{id} | 详情 |
| routes | PUT /api/routes/{id} | 更新 |
| routes | DELETE /api/routes/{id} | 删除 |
| routes | PATCH /api/routes/{id}/status | 状态流转 |
| health-checks | POST /api/health-checks | 创建HealthCheck |
| health-checks | GET /api/health-checks | 列表查询（分页+筛选） |
| health-checks | GET /api/health-checks/{id} | 详情 |
| health-checks | PUT /api/health-checks/{id} | 更新 |
| health-checks | DELETE /api/health-checks/{id} | 删除 |
| health-checks | PATCH /api/health-checks/{id}/status | 状态流转 |
| strategies | POST /api/strategies | 创建Strategy |
| strategies | GET /api/strategies | 列表查询（分页+筛选） |
| strategies | GET /api/strategies/{id} | 详情 |
| strategies | PUT /api/strategies/{id} | 更新 |
| strategies | DELETE /api/strategies/{id} | 删除 |
| strategies | PATCH /api/strategies/{id}/status | 状态流转 |
| access-logs | POST /api/access-logs | 创建AccessLog |
| access-logs | GET /api/access-logs | 列表查询（分页+筛选） |
| access-logs | GET /api/access-logs/{id} | 详情 |
| access-logs | DELETE /api/access-logs/{id} | 删除 |
| allow-lists | POST /api/allow-lists | 创建AllowList |
| allow-lists | GET /api/allow-lists | 列表查询（分页+筛选） |
| allow-lists | GET /api/allow-lists/{id} | 详情 |
| allow-lists | PUT /api/allow-lists/{id} | 更新 |
| allow-lists | DELETE /api/allow-lists/{id} | 删除 |
| allow-lists | PATCH /api/allow-lists/{id}/status | 状态流转 |
| stats | GET /api/stats/overview | 全局统计总览 |

## 说明

- 数据存储为内存实现，重启即清空。
- 金额/时间等数值字段统一采用整数（分 / 秒 / 毫秒 / Unix 时间戳），避免浮点精度问题。
