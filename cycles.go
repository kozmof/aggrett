package aggrep

import "time"

// CreateCycles creates repeated factors from start to end using a fixed interval.
func CreateCycles(
	start, end time.Time,
	step int,
	intervalType IntervalType,
	factor Factor,
	value float64,
	tag string,
	genID func() string,
) []SeqFactor {
	times := GenerateTimeSeries(start, end, step, intervalType)
	sequence := make([]SeqFactor, 0, len(times))

	for _, t := range times {
		sequence = append(sequence, SeqFactor{
			ID:     genID(),
			Tag:    tag,
			Value:  value,
			Factor: factor,
			Time:   t,
		})
	}

	return sequence
}
