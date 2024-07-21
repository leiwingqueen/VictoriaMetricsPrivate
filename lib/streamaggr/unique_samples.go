package streamaggr

import (
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/fasttime"
	"sync"
)

// uniqueSamplesAggrState calculates output=unique_samples, e.g. the number of unique sample values.
type uniqueSamplesAggrState struct {
	m sync.Map
}

type uniqueSamplesStateValue struct {
	mu sync.Mutex
	// implement
	mp      map[float64]struct{}
	deleted bool
}

func newUniqueSamplesAggrState() *uniqueSamplesAggrState {
	return &uniqueSamplesAggrState{}
}

func (as *uniqueSamplesAggrState) pushSamples(samples []pushSample) {
	// implement unique_samples pushSamples
	for _, sample := range samples {
		for {
			outputKey := getOutputKey(sample.key)
			var value *uniqueSamplesStateValue
			v, ok := as.m.Load(outputKey)
			if !ok {
				value = &uniqueSamplesStateValue{mp: make(map[float64]struct{}), deleted: false}
				actual, loaded := as.m.LoadOrStore(outputKey, value)
				if loaded {
					value = actual.(*uniqueSamplesStateValue)
				}
			} else {
				value = v.(*uniqueSamplesStateValue)
			}
			// update value
			value.mu.Lock()
			deleted := value.deleted
			if !deleted {
				value.mp[sample.value] = struct{}{}
			}
			value.mu.Unlock()
			if !deleted {
				break
			}
		}
	}
}

func (as *uniqueSamplesAggrState) flushState(ctx *flushCtx, resetState bool) {
	// implement unique_samples flushState
	timestamp := int64(fasttime.UnixTimestamp()) * 1_000
	as.m.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(*uniqueSamplesStateValue)
		if resetState {
			as.m.Delete(key)
		}
		v.mu.Lock()
		if resetState {
			v.deleted = true
		}
		v.mu.Unlock()
		ctx.appendSeries(k, "unique_samples", timestamp, float64(len(v.mp)))
		return true
	})
}
