package aggrett

import (
	"fmt"
	"sort"
	"time"
)

type timeGroup struct {
	Time    time.Time
	Factors []SeqFactor
}

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

// groupByTime copies sequence, sorts by time, and groups either exact timestamps or interval buckets.
func groupByTime(sequence []SeqFactor, grouping *timeGrouping) []timeGroup {
	if len(sequence) == 0 {
		return nil
	}
	sorted := make([]SeqFactor, len(sequence))
	copy(sorted, sequence)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Time.Before(sorted[j].Time)
	})

	var groups []timeGroup
	var current timeGroup
	for _, f := range sorted {
		groupTime := f.Time
		if grouping != nil {
			groupTime = bucketStart(f.Time, *grouping)
		}

		if len(current.Factors) == 0 {
			current = timeGroup{Time: groupTime, Factors: []SeqFactor{f}}
			continue
		}

		if current.Time.Equal(groupTime) {
			current.Factors = append(current.Factors, f)
			continue
		}

		groups = append(groups, current)
		current = timeGroup{Time: groupTime, Factors: []SeqFactor{f}}
	}
	return append(groups, current)
}

// AccumulateSequence sorts by time and accumulates values into time buckets.
func AccumulateSequence(sequence []SeqFactor, baseValue float64) []AccumCore {
	return accumulateSequence(sequence, baseValue, nil)
}

// AccumulateSequenceByInterval sorts by time and accumulates values into interval buckets.
// Result times are the start of each bucket.
func AccumulateSequenceByInterval(sequence []SeqFactor, baseValue float64, step int, intervalType IntervalType) []AccumCore {
	return accumulateSequence(sequence, baseValue, &timeGrouping{
		Step:         step,
		IntervalType: intervalType,
	})
}

func accumulateSequence(sequence []SeqFactor, baseValue float64, grouping *timeGrouping) []AccumCore {
	groups := groupByTime(sequence, grouping)
	if len(groups) == 0 {
		return []AccumCore{}
	}

	result := make([]AccumCore, 0, len(groups))
	store := baseValue
	for _, group := range groups {
		ids := make([]string, 0, len(group.Factors))
		for _, f := range group.Factors {
			ids = append(ids, f.ID)
			store = Accumulate(f.Factor, store, f.Value)
		}
		result = append(result, AccumCore{IDs: ids, Time: group.Time, Store: store})
	}
	return result
}
