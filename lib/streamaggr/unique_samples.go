package streamaggr

import (
	"sync"
)

// uniqueSamplesAggrState calculates output=unique_samples, e.g. the number of unique sample values.
type uniqueSamplesAggrState struct {
	m sync.Map
}

type uniqueSamplesStateValue struct {
	mu sync.Mutex
	// TODO: implement
	deleted bool
}

func newUniqueSamplesAggrState() *uniqueSamplesAggrState {
	return &uniqueSamplesAggrState{}
}

func (as *uniqueSamplesAggrState) pushSamples(samples []pushSample) {
	// TODO: implement unique_samples pushSamples
}

func (as *uniqueSamplesAggrState) flushState(ctx *flushCtx, resetState bool) {
	// TODO: implement unique_samples flushState
}
