package aggrep

import "sort"

func addBreakdown(breakdown Breakdown, factor SeqFactor) Breakdown {
	newBreakdown := cloneBreakdown(breakdown)
	existing, ok := newBreakdown[factor.Tag]
	if ok {
		ids := append(append([]string{}, existing.IDs...), factor.ID)
		newBreakdown[factor.Tag] = BreakdownEntry{
			IDs:   ids,
			Store: Accumulate(factor.Factor, existing.Store, factor.Value),
		}
	} else {
		newBreakdown[factor.Tag] = BreakdownEntry{
			IDs:   []string{factor.ID},
			Store: Accumulate(factor.Factor, 0, factor.Value),
		}
	}
	return newBreakdown
}

// Aggregate groups factors by time and returns running totals with per-time breakdown.
// If filter is empty, all tags are included. Otherwise only listed tags are included.
func Aggregate(sequence []SeqFactor, baseValue float64, filter []string) []Accum {
	if len(sequence) < 1 {
		return []Accum{}
	}

	sorted := make([]SeqFactor, len(sequence))
	copy(sorted, sequence)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Time.Before(sorted[j].Time)
	})

	filtered := sorted
	if len(filter) > 0 {
		filterSet := make(map[string]struct{}, len(filter))
		for _, tag := range filter {
			filterSet[tag] = struct{}{}
		}
		filtered = make([]SeqFactor, 0, len(sorted))
		for _, f := range sorted {
			if _, ok := filterSet[f.Tag]; ok {
				filtered = append(filtered, f)
			}
		}
	}

	if len(filtered) < 1 {
		return []Accum{}
	}

	firstFactor := filtered[0]
	restFactors := filtered[1:]

	accums := make([]Accum, 0, len(filtered))
	timePos := firstFactor.Time
	accum := Accum{
		IDs:   []string{firstFactor.ID},
		Tags:  []string{firstFactor.Tag},
		Time:  timePos,
		Store: Accumulate(firstFactor.Factor, baseValue, firstFactor.Value),
		Breakdown: Breakdown{
			firstFactor.Tag: {
				IDs:   []string{firstFactor.ID},
				Store: Accumulate(firstFactor.Factor, 0, firstFactor.Value),
			},
		},
	}

	for _, seqFactor := range restFactors {
		prevValue := accum.Store

		if timePos.Equal(seqFactor.Time) {
			accum = Accum{
				IDs:       append(append([]string{}, accum.IDs...), seqFactor.ID),
				Tags:      appendTagIfMissing(accum.Tags, seqFactor.Tag),
				Time:      timePos,
				Store:     Accumulate(seqFactor.Factor, prevValue, seqFactor.Value),
				Breakdown: addBreakdown(accum.Breakdown, seqFactor),
			}
		} else {
			accums = append(accums, accum)
			newTimePos := seqFactor.Time
			accum = Accum{
				IDs:   []string{seqFactor.ID},
				Tags:  []string{seqFactor.Tag},
				Time:  newTimePos,
				Store: Accumulate(seqFactor.Factor, prevValue, seqFactor.Value),
				Breakdown: Breakdown{
					seqFactor.Tag: {
						IDs:   []string{seqFactor.ID},
						Store: Accumulate(seqFactor.Factor, 0, seqFactor.Value),
					},
				},
			}
			timePos = newTimePos
		}
	}

	accums = append(accums, accum)
	return accums
}

func appendTagIfMissing(tags []string, tag string) []string {
	for _, t := range tags {
		if t == tag {
			return append([]string{}, tags...)
		}
	}
	return append(append([]string{}, tags...), tag)
}

func cloneBreakdown(breakdown Breakdown) Breakdown {
	cloned := make(Breakdown, len(breakdown))
	for k, v := range breakdown {
		ids := append([]string{}, v.IDs...)
		cloned[k] = BreakdownEntry{Store: v.Store, IDs: ids}
	}
	return cloned
}
