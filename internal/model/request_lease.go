package model

// RequestLease 是从对象池借出的请求上下文载体，承载当前请求的
// 租户与链路追踪信息。由于对象会被跨请求复用，任何需要脱离池
// 生命周期的读取（例如异步日志）都必须先取 Snapshot，再归还对象，
// 以免读到下一个租户覆写后的字段。
type RequestLease struct {
	RequestID string
	Tenant    string
	Headers   map[string]string
}

// Snapshot 返回 lease 的值拷贝，并深拷贝 Headers map。
// 调用方应在归还对象之前获取快照，用于后续异步处理。
func (l *RequestLease) Snapshot() RequestLease {
	headers := make(map[string]string, len(l.Headers))
	for k, v := range l.Headers {
		headers[k] = v
	}
	return RequestLease{
		RequestID: l.RequestID,
		Tenant:    l.Tenant,
		Headers:   headers,
	}
}

// Reset 清空自身字段，供对象池在归还/复用前调用，确保上一个请求
// 的租户与 trace 不会残留到下一个请求。
func (l *RequestLease) Reset() {
	l.RequestID = ""
	l.Tenant = ""
	for k := range l.Headers {
		delete(l.Headers, k)
	}
}
