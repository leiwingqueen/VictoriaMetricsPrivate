package streamaggr

import (
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/fasttime"
	"sync"
)

// countSamplesAggrState calculates output=count_samples, e.g. the count of input samples.
type countSamplesAggrState struct {
	m sync.Map
}

type countSamplesStateValue struct {
	mu sync.Mutex
	// implement
	count   int64
	deleted bool
}

func newCountSamplesAggrState() *countSamplesAggrState {
	return &countSamplesAggrState{}
}

func (as *countSamplesAggrState) pushSamples(samples []pushSample) {
	// implement
	for _, sample := range samples {
		for {
			outputKey := getOutputKey(sample.key)
			var value *countSamplesStateValue
			v, ok := as.m.Load(outputKey)
			if !ok {
				value = &countSamplesStateValue{count: 0, deleted: false}
				actual, loaded := as.m.LoadOrStore(outputKey, value)
				if loaded {
					value = actual.(*countSamplesStateValue)
				}
			} else {
				value = v.(*countSamplesStateValue)
			}
			// update value
			value.mu.Lock()
			deleted := value.deleted
			if !deleted {
				value.count++
				value.count++
			}
			value.mu.Unlock()
			if !deleted {
				break
			}
		}
	}
}

func (as *countSamplesAggrState) flushState(ctx *flushCtx, resetState bool) {
	// implement
	timestamp := int64(fasttime.UnixTimestamp()) * 1_000
	as.m.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(*countSamplesStateValue)
		if resetState {
			as.m.Delete(key)
		}
		avg := float64(v.count)
		v.mu.Lock()
		if resetState {
			v.deleted = true
		}
		v.mu.Unlock()
		ctx.appendSeries(k, "count_samples", timestamp, avg)
		return true
	})
}
