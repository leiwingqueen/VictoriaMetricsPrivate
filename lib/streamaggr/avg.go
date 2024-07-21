package streamaggr

import (
	"sync"
)

// avgAggrState calculates output=avg, e.g. the average value over input samples.
type avgAggrState struct {
	m sync.Map
}

type avgStateValue struct {
	mu      sync.Mutex
	sum     float64
	count   int64
	deleted bool
}

func newAvgAggrState() *avgAggrState {
	return &avgAggrState{}
}

func (as *avgAggrState) pushSamples(samples []pushSample) {
	// TODO: implement pushSamples
	// hint:
	// - get the output key from the sample key
	// - get the value from the map
	// - if the value is not present, create a new value and store it in the map
	// - remember to lock the value before updating it
	// - if the value is deleted, try to obtain and update it again
}

func (as *avgAggrState) flushState(ctx *flushCtx, resetState bool) {
	// TODO: implement flushState
	// hint:
	// - if resetState is true, then delete the entry from the map
	// - get the value as avgStateValue and calculate the avg
	// - if resetState is true, update the deleted flag
	// - remember to append the series to the ctx
}
