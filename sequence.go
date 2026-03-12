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
func Accumulate(factor Factor, prevValue, value float64) (float64, error) {
	switch factor {
	case FactorPlus:
		return prevValue + value, nil
	case FactorMinus:
		return prevValue - value, nil
	default:
		return 0, fmt.Errorf("unknown factor type: %q", factor)
	}
}

// groupByTime copies sequence, sorts by time, and groups either exact timestamps or interval buckets.
func groupByTime(sequence []SeqFactor, grouping *timeGrouping) ([]timeGroup, error) {
	if len(sequence) == 0 {
		return nil, nil
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
			var err error
			groupTime, err = bucketStart(f.Time, *grouping)
			if err != nil {
				return nil, err
			}
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
	return append(groups, current), nil
}

// AccumulateSequence sorts by time and accumulates values into time buckets.
func AccumulateSequence(sequence []SeqFactor, baseValue float64) ([]AccumCore, error) {
	return accumulateSequence(sequence, baseValue, nil)
}

// AccumulateSequenceByInterval sorts by time and accumulates values into interval buckets.
// Result times are the start of each bucket.
func AccumulateSequenceByInterval(sequence []SeqFactor, baseValue float64, step int, intervalType IntervalType) ([]AccumCore, error) {
	return accumulateSequence(sequence, baseValue, &timeGrouping{
		Step:         step,
		IntervalType: intervalType,
	})
}

func accumulateSequence(sequence []SeqFactor, baseValue float64, grouping *timeGrouping) ([]AccumCore, error) {
	groups, err := groupByTime(sequence, grouping)
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return []AccumCore{}, nil
	}

	result := make([]AccumCore, 0, len(groups))
	store := baseValue
	for _, group := range groups {
		ids := make([]string, 0, len(group.Factors))
		for _, f := range group.Factors {
			ids = append(ids, f.ID)
			store, err = Accumulate(f.Factor, store, f.Value)
			if err != nil {
				return nil, err
			}
		}
		result = append(result, AccumCore{IDs: ids, Time: group.Time, Store: store})
	}
	return result, nil
}
