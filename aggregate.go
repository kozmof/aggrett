package aggrett

// Aggregate groups factors by exact time and returns running totals with per-time breakdown.
// If filter is empty, all tags are included. Otherwise only listed tags are included.
func Aggregate(sequence []SeqFactor, baseValue float64, filter []string) ([]Accum, error) {
	return aggregate(sequence, baseValue, filter, nil)
}

// AggregateByInterval groups factors into interval buckets and returns running totals with per-bucket breakdown.
// Result times are the start of each bucket.
func AggregateByInterval(sequence []SeqFactor, baseValue float64, filter []string, step int, intervalType IntervalType) ([]Accum, error) {
	return aggregate(sequence, baseValue, filter, &TimeGrouping{
		Step:         step,
		IntervalType: intervalType,
	})
}

// AggregateByGrouping is like AggregateByInterval but accepts a pre-built TimeGrouping,
// allowing callers to construct and validate it once and reuse.
func AggregateByGrouping(sequence []SeqFactor, baseValue float64, filter []string, grouping TimeGrouping) ([]Accum, error) {
	return aggregate(sequence, baseValue, filter, &grouping)
}

func aggregate(sequence []SeqFactor, baseValue float64, filter []string, grouping *TimeGrouping) ([]Accum, error) {
	filtered := sequence
	if len(filter) > 0 {
		filtered = FilterByTag(sequence, filter)
	}

	groups, err := groupByTime(filtered, grouping)
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return []Accum{}, nil
	}

	result := make([]Accum, 0, len(groups))
	store := baseValue
	for _, group := range groups {
		ids := make([]string, 0, len(group.Factors))
		tags := make([]string, 0, len(group.Factors))
		tagsSeen := make(map[string]struct{})
		breakdown := make(Breakdown)

		for _, f := range group.Factors {
			ids = append(ids, f.ID)
			store, err = Accumulate(f.Factor, store, f.Value)
			if err != nil {
				return nil, err
			}

			if _, ok := tagsSeen[f.Tag]; !ok {
				tagsSeen[f.Tag] = struct{}{}
				tags = append(tags, f.Tag)
			}

			entry := breakdown[f.Tag]
			delta, err := Accumulate(f.Factor, entry.Delta, f.Value)
			if err != nil {
				return nil, err
			}
			breakdown[f.Tag] = BreakdownEntry{
				IDs:   append(entry.IDs, f.ID),
				Delta: delta,
			}
		}

		result = append(result, Accum{
			AccumCore: AccumCore{IDs: ids, Time: group.Time, Store: store},
			Tags:      tags,
			Breakdown: breakdown,
		})
	}
	return result, nil
}
