package service

import (
	"reverseproxy/internal/model"
	"reverseproxy/internal/store"
)

type RequestLeaseProcessor struct{ pool *store.RequestLeasePool }

func NewRequestLeaseProcessor(pool *store.RequestLeasePool) *RequestLeaseProcessor {
	return &RequestLeaseProcessor{pool: pool}
}

// fillLease 用当前请求的字段填充借出的 lease。
// 出于对象复用的隔离考虑，这里先清空再写入，确保即便池未在
// Acquire 时 reset，也不会残留上一个请求的 tenant/trace。
func fillLease(lease *model.RequestLease, requestID, tenant string, headers map[string]string) {
	lease.RequestID = requestID
	lease.Tenant = tenant
	for key := range lease.Headers {
		delete(lease.Headers, key)
	}
	for key, value := range headers {
		lease.Headers[key] = value
	}
}

// StartAsync 借出 lease 填充当前请求上下文后，在归还前先取快照，
// 再把对象归还入池；异步 goroutine 只读快照值，绝不触碰归还后
// 可能已被下一个租户复用的共享对象。这样租户与 trace 不会串请求。
func (p *RequestLeaseProcessor) StartAsync(requestID, tenant string, headers map[string]string, release <-chan struct{}) <-chan string {
	lease := p.pool.Acquire()
	fillLease(lease, requestID, tenant, headers)
	snap := lease.Snapshot()
	p.pool.Release(lease)
	result := make(chan string, 1)
	go func() {
		<-release
		result <- snap.Tenant + ":" + snap.Headers["trace"]
		close(result)
	}()
	return result
}

func (p *RequestLeaseProcessor) ProcessSync(requestID, tenant string, headers map[string]string) string {
	lease := p.pool.Acquire()
	fillLease(lease, requestID, tenant, headers)
	value := lease.Tenant + ":" + lease.Headers["trace"]
	p.pool.Release(lease)
	return value
}
