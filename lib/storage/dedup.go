package storage

import (
	"time"
)

// SetDedupInterval sets the deduplication interval, which is applied to raw samples during data ingestion and querying.
//
// De-duplication is disabled if dedupInterval is 0.
//
// This function must be called before initializing the storage.
func SetDedupInterval(dedupInterval time.Duration) {
	globalDedupInterval = dedupInterval.Milliseconds()
}

// GetDedupInterval returns the dedup interval in milliseconds, which has been set via SetDedupInterval.
func GetDedupInterval() int64 {
	return globalDedupInterval
}

var globalDedupInterval int64

func isDedupEnabled() bool {
	return globalDedupInterval > 0
}

// DeduplicateSamples removes samples from src* if they are closer to each other than dedupInterval in milliseconds.
// DeduplicateSamples treats StaleNaN (Prometheus stale markers) as values and doesn't skip them on purpose - see
// https://github.com/VictoriaMetrics/VictoriaMetrics/issues/5587
func DeduplicateSamples(srcTimestamps []int64, srcValues []float64, dedupInterval int64) ([]int64, []float64) {
	if !needsDedup(srcTimestamps, dedupInterval) {
		// Fast path - nothing to deduplicate
		return srcTimestamps, srcValues
	}
	nextTs := srcTimestamps[0] - srcTimestamps[0]%dedupInterval
	// tsNext -= tsNext % dedupInterval
	// dstTimestamps := srcTimestamps[:0]
	// dstValues := srcValues[:0]
	n := len(srcTimestamps)
	p := 0
	for i := 0; i < n; i++ {
		ts := srcTimestamps[i]
		if ts < nextTs {
			continue
		}
		value := srcValues[i]
		// Choose the maximum value with the timestamp equal to tsPrev.
		// See https://github.com/VictoriaMetrics/VictoriaMetrics/issues/3333
		j := i
		for ; j < n && srcTimestamps[j] == ts; j++ {
			value = max(value, srcValues[j])
		}
		i = j - 1
		srcTimestamps[p] = ts
		srcValues[p] = value
		p++
		// move to next time window
		nextTs += dedupInterval
		if nextTs < ts {
			nextTs = ts + dedupInterval
			nextTs -= nextTs % dedupInterval
		}
	}
	return srcTimestamps[:p], srcValues[:p]
}

func deduplicateSamplesDuringMerge(srcTimestamps, srcValues []int64, dedupInterval int64) ([]int64, []int64) {
	if !needsDedup(srcTimestamps, dedupInterval) {
		// Fast path - nothing to deduplicate
		return srcTimestamps, srcValues
	}
	// implement deduplicate samples
	// don't need to allocate new memory to save the
	// [0,dedupInterval) just save one metric
	nextTs := srcTimestamps[0] - srcTimestamps[0]%dedupInterval
	n := len(srcTimestamps)
	// pointer of next write
	p := 0
	for i := 0; i < n; i++ {
		value := srcValues[i]
		ts := srcTimestamps[i]
		if ts < nextTs {
			continue
		}
		// 相同的ts取最大值
		j := i
		for ; j < n && srcTimestamps[i] == srcTimestamps[j]; j++ {
			value = max(value, srcValues[j])
		}
		i = j - 1
		srcTimestamps[p] = ts
		srcValues[p] = value
		p++
		nextTs = ts + dedupInterval
		nextTs -= nextTs % dedupInterval
	}
	return srcTimestamps[:p], srcValues[:p]
}

func needsDedup(timestamps []int64, dedupInterval int64) bool {
	if len(timestamps) < 2 || dedupInterval <= 0 {
		return false
	}
	tsNext := timestamps[0] + dedupInterval - 1
	tsNext -= tsNext % dedupInterval
	for _, ts := range timestamps[1:] {
		if ts <= tsNext {
			return true
		}
		tsNext += dedupInterval
		if tsNext < ts {
			tsNext = ts + dedupInterval - 1
			tsNext -= tsNext % dedupInterval
		}
	}
	return false
}
