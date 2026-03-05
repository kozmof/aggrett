package aggrep

import (
	"fmt"
	"sort"
)

// Accumulate applies a factor to a previous value.
func Accumulate(factor Factor, prevValue, value float64) float64 {
	switch factor {
	case FactorPlus:
		return prevValue + value
	case FactorMinus:
		return prevValue - value
	default:
		panic(fmt.Sprintf("unknown factor type: %q", factor))
	}
}

// groupByTime copies sequence, sorts by time, and returns slices of same-time factors.
func groupByTime(sequence []SeqFactor) [][]SeqFactor {
	if len(sequence) == 0 {
		return nil
	}
	sorted := make([]SeqFactor, len(sequence))
	copy(sorted, sequence)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Time.Before(sorted[j].Time)
	})

	var groups [][]SeqFactor
	var current []SeqFactor
	for _, f := range sorted {
		if len(current) == 0 || current[0].Time.Equal(f.Time) {
			current = append(current, f)
		} else {
			groups = append(groups, current)
			current = []SeqFactor{f}
		}
	}
	return append(groups, current)
}

// AccumulateSequence sorts by time and accumulates values into time buckets.
func AccumulateSequence(sequence []SeqFactor, baseValue float64) []AccumCore {
	groups := groupByTime(sequence)
	if len(groups) == 0 {
		return []AccumCore{}
	}

	result := make([]AccumCore, 0, len(groups))
	store := baseValue
	for _, group := range groups {
		ids := make([]string, 0, len(group))
		for _, f := range group {
			ids = append(ids, f.ID)
			store = Accumulate(f.Factor, store, f.Value)
		}
		result = append(result, AccumCore{IDs: ids, Time: group[0].Time, Store: store})
	}
	return result
}
