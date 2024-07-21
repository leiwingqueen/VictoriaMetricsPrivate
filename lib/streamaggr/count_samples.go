package streamaggr

import (
	"sync"
)

// countSamplesAggrState calculates output=count_samples, e.g. the count of input samples.
type countSamplesAggrState struct {
	m sync.Map
}

type countSamplesStateValue struct {
	mu sync.Mutex
	// TODO: implement
	deleted bool
}

func newCountSamplesAggrState() *countSamplesAggrState {
	return &countSamplesAggrState{}
}

func (as *countSamplesAggrState) pushSamples(samples []pushSample) {
	// TODO: implement
}

func (as *countSamplesAggrState) flushState(ctx *flushCtx, resetState bool) {
	// TODO: implement
}
