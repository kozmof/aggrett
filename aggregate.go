package aggrep

import "sort"

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

// Aggregate groups factors by time and returns running totals with per-time breakdown.
// If filter is empty, all tags are included. Otherwise only listed tags are included.
func Aggregate(sequence []SeqFactor, baseValue float64, filter []string) []Accum {
	filtered := sequence
	if len(filter) > 0 {
		filtered = FilterByTag(sequence, filter)
	}

	groups := groupByTime(filtered)
	if len(groups) == 0 {
		return []Accum{}
	}

	result := make([]Accum, 0, len(groups))
	store := baseValue
	for _, group := range groups {
		ids := make([]string, 0, len(group))
		tags := make([]string, 0)
		tagsSeen := make(map[string]struct{})
		breakdown := make(Breakdown)

		for _, f := range group {
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
			Time:      group[0].Time,
			Store:     store,
			Breakdown: breakdown,
		})
	}
	return result
}
