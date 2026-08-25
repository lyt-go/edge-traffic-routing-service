package service

import (
	"context"

	"reverseproxy/internal/store"
)

type ContextProbeRunner struct{ probes *store.ContextProbeStore }

func NewContextProbeRunner(probes *store.ContextProbeStore) *ContextProbeRunner {
	return &ContextProbeRunner{probes: probes}
}

func (r *ContextProbeRunner) Run(ctx context.Context, upstreamID string) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		if err = r.probes.Check(ctx, upstreamID); err == nil {
			return nil
		}
	}
	return err
}
