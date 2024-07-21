package streamaggr

import (
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/fasttime"
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
	// implement pushSamples
	// hint:
	// - get the output key from the sample key
	// - get the value from the map
	// - if the value is not present, create a new value and store it in the map
	// - remember to lock the value before updating it
	// - if the value is deleted, try to obtain and update it again
	for _, sample := range samples {
		for {
			outputKey := getOutputKey(sample.key)
			var value *avgStateValue
			v, ok := as.m.Load(outputKey)
			if !ok {
				value = &avgStateValue{sum: 0, count: 0, deleted: false}
				actual, loaded := as.m.LoadOrStore(outputKey, value)
				if loaded {
					value = actual.(*avgStateValue)
				}
			} else {
				value = v.(*avgStateValue)
			}
			// update value
			value.mu.Lock()
			deleted := value.deleted
			if !deleted {
				value.sum += sample.value
				value.count++
			}
			value.mu.Unlock()
			if !deleted {
				break
			}
		}
	}
}

func (as *avgAggrState) flushState(ctx *flushCtx, resetState bool) {
	// implement flushState
	// hint:
	// - if resetState is true, then delete the entry from the map
	// - get the value as avgStateValue and calculate the avg
	// - if resetState is true, update the deleted flag
	// - remember to append the series to the ctx
	timestamp := int64(fasttime.UnixTimestamp()) * 1_000
	as.m.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(*avgStateValue)
		if resetState {
			as.m.Delete(key)
		}
		avg := v.sum / float64(v.count)
		v.mu.Lock()
		if resetState {
			v.deleted = true
		}
		v.mu.Unlock()
		ctx.appendSeries(k, "avg", timestamp, avg)
		return true
	})
}
