package streamaggr

import (
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/fasttime"
	"sync"
)

// sumSamplesAggrState calculates output=sum_samples, e.g. the sum over input samples.
type sumSamplesAggrState struct {
	m sync.Map
}

type sumSamplesStateValue struct {
	mu sync.Mutex
	// implement
	sum     float64
	deleted bool
}

func newSumSamplesAggrState() *sumSamplesAggrState {
	return &sumSamplesAggrState{}
}

func (as *sumSamplesAggrState) pushSamples(samples []pushSample) {
	// implement
	for _, sample := range samples {
		for {
			outputKey := getOutputKey(sample.key)
			var value *sumSamplesStateValue
			v, ok := as.m.Load(outputKey)
			if !ok {
				value = &sumSamplesStateValue{sum: 0, deleted: false}
				actual, loaded := as.m.LoadOrStore(outputKey, value)
				if loaded {
					value = actual.(*sumSamplesStateValue)
				}
			} else {
				value = v.(*sumSamplesStateValue)
			}
			// update value
			value.mu.Lock()
			deleted := value.deleted
			if !deleted {
				value.sum += sample.value
			}
			value.mu.Unlock()
			if !deleted {
				break
			}
		}
	}
}

func (as *sumSamplesAggrState) flushState(ctx *flushCtx, resetState bool) {
	// implement
	timestamp := int64(fasttime.UnixTimestamp()) * 1_000
	as.m.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(*sumSamplesStateValue)
		if resetState {
			as.m.Delete(key)
		}
		v.mu.Lock()
		if resetState {
			v.deleted = true
		}
		v.mu.Unlock()
		ctx.appendSeries(k, "sum_samples", timestamp, v.sum)
		return true
	})
}
