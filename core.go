package aggrep

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
