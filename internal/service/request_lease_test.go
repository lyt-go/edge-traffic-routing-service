package service_test

import (
	"testing"

	"reverseproxy/internal/service"
	"reverseproxy/internal/store"
)

func TestPooledRequestStaysOwnedUntilAsyncLogCompletes(t *testing.T) {
	processor := service.NewRequestLeaseProcessor(store.NewRequestLeasePool())
	release := make(chan struct{})
	firstLog := processor.StartAsync("req-a", "tenant-alpha", map[string]string{"trace": "trace-a"}, release)
	if got := processor.ProcessSync("req-b", "tenant-beta", nil); got != "tenant-beta:" {
		t.Fatalf("second request inherited pooled fields: got=%q", got)
	}
	close(release)
	if got := <-firstLog; got != "tenant-alpha:trace-a" {
		t.Fatalf("first async log used the next request identity: got=%q", got)
	}
	if got := processor.ProcessSync("req-c", "", nil); got != ":" {
		t.Fatalf("later empty request retained an old tenant or trace: got=%q", got)
	}
}
