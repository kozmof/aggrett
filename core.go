package aggrep

import "sort"

// AccumulateSequence sorts by time and accumulates values into time buckets.
func AccumulateSequence(sequence []SeqFactor, baseValue float64) []AccumCore {
	if len(sequence) < 1 {
		return []AccumCore{}
	}

	sorted := make([]SeqFactor, len(sequence))
	copy(sorted, sequence)
	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].Time.Before(sorted[j].Time)
	})

	firstFactor := sorted[0]
	restFactors := sorted[1:]

	accums := make([]AccumCore, 0, len(sorted))
	timePos := firstFactor.Time
	accum := AccumCore{
		IDs:   []string{firstFactor.ID},
		Time:  timePos,
		Store: Accumulate(firstFactor.Factor, baseValue, firstFactor.Value),
	}

	for _, seqFactor := range restFactors {
		prevValue := accum.Store

		if timePos.Equal(seqFactor.Time) {
			ids := append(append([]string{}, accum.IDs...), seqFactor.ID)
			accum = AccumCore{
				IDs:   ids,
				Time:  timePos,
				Store: Accumulate(seqFactor.Factor, prevValue, seqFactor.Value),
			}
		} else {
			accums = append(accums, accum)
			timePos = seqFactor.Time
			accum = AccumCore{
				IDs:   []string{seqFactor.ID},
				Time:  timePos,
				Store: Accumulate(seqFactor.Factor, prevValue, seqFactor.Value),
			}
		}
	}

	accums = append(accums, accum)
	return accums
}
