package streamaggr

import (
	"sync"
)

// sumSamplesAggrState calculates output=sum_samples, e.g. the sum over input samples.
type sumSamplesAggrState struct {
	m sync.Map
}

type sumSamplesStateValue struct {
	mu sync.Mutex
	// implement
	deleted bool
}

func newSumSamplesAggrState() *sumSamplesAggrState {
	return &sumSamplesAggrState{}
}

func (as *sumSamplesAggrState) pushSamples(samples []pushSample) {
	// TODO: implement
}

func (as *sumSamplesAggrState) flushState(ctx *flushCtx, resetState bool) {
	// TODO: implement
}
