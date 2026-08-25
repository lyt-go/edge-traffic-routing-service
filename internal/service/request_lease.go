package service

import (
	"reverseproxy/internal/model"
	"reverseproxy/internal/store"
)

type RequestLeaseProcessor struct{ pool *store.RequestLeasePool }

func NewRequestLeaseProcessor(pool *store.RequestLeasePool) *RequestLeaseProcessor {
	return &RequestLeaseProcessor{pool: pool}
}

func fillLease(lease *model.RequestLease, requestID, tenant string, headers map[string]string) {
	lease.RequestID = requestID
	if tenant != "" {
		lease.Tenant = tenant
	}
	for key, value := range headers {
		lease.Headers[key] = value
	}
}

func (p *RequestLeaseProcessor) StartAsync(requestID, tenant string, headers map[string]string, release <-chan struct{}) <-chan string {
	lease := p.pool.Acquire()
	fillLease(lease, requestID, tenant, headers)
	result := make(chan string, 1)
	go func() {
		<-release
		result <- lease.Tenant + ":" + lease.Headers["trace"]
		close(result)
	}()
	p.pool.Release(lease)
	return result
}

func (p *RequestLeaseProcessor) ProcessSync(requestID, tenant string, headers map[string]string) string {
	lease := p.pool.Acquire()
	fillLease(lease, requestID, tenant, headers)
	value := lease.Tenant + ":" + lease.Headers["trace"]
	p.pool.Release(lease)
	return value
}
