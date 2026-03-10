package aggrett

// Aggregate groups factors by exact time and returns running totals with per-time breakdown.
// If filter is empty, all tags are included. Otherwise only listed tags are included.
func Aggregate(sequence []SeqFactor, baseValue float64, filter []string) []Accum {
	return aggregate(sequence, baseValue, filter, nil)
}

// AggregateByInterval groups factors into interval buckets and returns running totals with per-bucket breakdown.
// Result times are the start of each bucket.
func AggregateByInterval(sequence []SeqFactor, baseValue float64, filter []string, step int, intervalType IntervalType) []Accum {
	return aggregate(sequence, baseValue, filter, &timeGrouping{
		Step:         step,
		IntervalType: intervalType,
	})
}

func aggregate(sequence []SeqFactor, baseValue float64, filter []string, grouping *timeGrouping) []Accum {
	filtered := sequence
	if len(filter) > 0 {
		filtered = FilterByTag(sequence, filter)
	}

	groups := groupByTime(filtered, grouping)
	if len(groups) == 0 {
		return []Accum{}
	}

	result := make([]Accum, 0, len(groups))
	store := baseValue
	for _, group := range groups {
		ids := make([]string, 0, len(group.Factors))
		tags := make([]string, 0)
		tagsSeen := make(map[string]struct{})
		breakdown := make(Breakdown)

		for _, f := range group.Factors {
			ids = append(ids, f.ID)
			store = Accumulate(f.Factor, store, f.Value)

			if _, ok := tagsSeen[f.Tag]; !ok {
				tagsSeen[f.Tag] = struct{}{}
				tags = append(tags, f.Tag)
			}

			entry := breakdown[f.Tag]
			breakdown[f.Tag] = BreakdownEntry{
				IDs:   append(entry.IDs, f.ID),
				Delta: Accumulate(f.Factor, entry.Delta, f.Value),
			}
		}

		result = append(result, Accum{
			IDs:       ids,
			Tags:      tags,
			Time:      group.Time,
			Store:     store,
			Breakdown: breakdown,
		})
	}
	return result
}
