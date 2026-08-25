package service

import (
	"context"
	"sync"

	"reverseproxy/internal/model"
	"reverseproxy/internal/store"
)

type ProbeBatchCollector struct{ stream *store.ProbeStream }

func NewProbeBatchCollector(stream *store.ProbeStream) *ProbeBatchCollector {
	return &ProbeBatchCollector{stream: stream}
}

func (c *ProbeBatchCollector) Collect(ctx context.Context, upstreamIDs []string) ([]model.ProbeBatchResult, error) {
	input, _ := c.stream.Start(ctx, upstreamIDs)
	processed := make(chan model.ProbeBatchResult)
	go func() {
		var workers sync.WaitGroup
		for result := range input {
			go func(item model.ProbeBatchResult) {
				workers.Add(1)
				defer workers.Done()
				processed <- item
			}(result)
		}
		workers.Wait()
		close(processed)
	}()
	collected := make([]model.ProbeBatchResult, 0, len(upstreamIDs))
	for result := range processed {
		collected = append(collected, result)
	}
	return collected, nil
}
