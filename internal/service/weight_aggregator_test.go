package service_test

import (
	"testing"

	"reverseproxy/internal/model"
	"reverseproxy/internal/service"
	"reverseproxy/internal/store"
)

func TestWeightSnapshotSurvivesConcurrentRefresh(t *testing.T) {
	st := store.NewRouteWeightStore()
	st.Put(model.RouteWeight{UpstreamID: "edge-a", Weight: 1})
	aggregator := service.NewWeightAggregator(st)
	release := make(chan struct{})
	total := aggregator.Begin(release)
	st.Put(model.RouteWeight{UpstreamID: "edge-a", Weight: 9})
	st.Put(model.RouteWeight{UpstreamID: "edge-b", Weight: 4})
	close(release)

	if got := <-total; got != 1 {
		t.Fatalf("in-flight route total changed during refresh: got=%d want=1", got)
	}
	last := aggregator.LastSnapshot()
	if len(last) != 1 || last["edge-a"].Weight != 1 {
		t.Fatalf("saved route snapshot was polluted: %#v", last)
	}
	last["edge-a"].Weight = 99
	if got := aggregator.LastSnapshot()["edge-a"].Weight; got != 1 {
		t.Fatalf("caller mutated the aggregator cache: got=%d", got)
	}
	raw := st.Snapshot()
	raw["edge-a"].Weight = 77
	if got := st.Snapshot()["edge-a"].Weight; got != 9 {
		t.Fatalf("store snapshot exposed live route weights: got=%d", got)
	}
}
