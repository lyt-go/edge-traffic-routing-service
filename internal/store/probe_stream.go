package store

import (
	"context"
	"fmt"

	"reverseproxy/internal/model"
)

type ProbeStream struct{}

func NewProbeStream() *ProbeStream { return &ProbeStream{} }

func (s *ProbeStream) Start(ctx context.Context, upstreamIDs []string) (<-chan model.ProbeBatchResult, <-chan error) {
	results := make(chan model.ProbeBatchResult)
	errs := make(chan error)
	go func() {
		for _, id := range upstreamIDs {
			if id == "unreachable" {
				errs <- fmt.Errorf("probe %s failed", id)
				return
			}
			select {
			case results <- model.ProbeBatchResult{UpstreamID: id, Reachable: true}:
			case <-ctx.Done():
				return
			}
		}
		close(results)
		close(errs)
	}()
	return results, errs
}
