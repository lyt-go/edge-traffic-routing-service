package model

type RequestLease struct {
	RequestID string
	Tenant    string
	Headers   map[string]string
}

func (l *RequestLease) Snapshot() RequestLease {
	return *l
}
