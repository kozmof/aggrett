package aggrett

import (
	"fmt"
	"time"
)

// IntervalType defines supported interval units.
type IntervalType string

const (
	IntervalSeconds IntervalType = "seconds"
	IntervalMinutes IntervalType = "minutes"
	IntervalHours   IntervalType = "hours"
	IntervalDays    IntervalType = "days"
	IntervalWeeks   IntervalType = "weeks"
	IntervalMonths  IntervalType = "months"
	IntervalYears   IntervalType = "years"
)

// String implements fmt.Stringer.
func (t IntervalType) String() string { return string(t) }

// IsValid reports whether t is a known IntervalType.
func (t IntervalType) IsValid() bool {
	switch t {
	case IntervalSeconds, IntervalMinutes, IntervalHours,
		IntervalDays, IntervalWeeks, IntervalMonths, IntervalYears:
		return true
	}
	return false
}

// TimeGrouping describes an interval step and unit for bucket-based operations.
// Result times are the start of each bucket.
type TimeGrouping struct {
	Step         int
	IntervalType IntervalType
}

func (g TimeGrouping) validate() error {
	if g.Step <= 0 {
		return fmt.Errorf("grouping step must be positive: %d", g.Step)
	}
	if !g.IntervalType.IsValid() {
		return fmt.Errorf("unknown interval type: %q", g.IntervalType)
	}
	return nil
}

// AddInterval adds step units of intervalType to date.
//
// For IntervalMonths and IntervalYears, overflow is clamped to the last day of the
// target month using JS-like semantics: if the source day does not exist in the
// target month, the result is the last valid day of that month. The clamped day
// becomes the anchor for subsequent steps, so monthly cycles starting on the 31st
// will drift permanently (e.g. Jan-31 → Feb-29 → Mar-29, not Mar-31).
func AddInterval(date time.Time, step int, intervalType IntervalType) time.Time {
	switch intervalType {
	case IntervalSeconds:
		return date.Add(time.Duration(step) * time.Second)
	case IntervalMinutes:
		return date.Add(time.Duration(step) * time.Minute)
	case IntervalHours:
		return date.Add(time.Duration(step) * time.Hour)
	case IntervalDays:
		return date.AddDate(0, 0, step)
	case IntervalWeeks:
		return date.AddDate(0, 0, step*7)
	case IntervalMonths:
		return addMonthsWithOverflowClamp(date, step)
	case IntervalYears:
		return addYearsWithOverflowClamp(date, step)
	default:
		panic("unknown interval type")
	}
}

// bucketStart returns the start of the bucket that date falls into for the given grouping.
//
// IntervalDays and IntervalWeeks use epoch-aligned buckets (anchored at 1970-01-01 UTC
// in the date's location) so that bucket boundaries are consistent across month and
// year boundaries regardless of step size.
func bucketStart(date time.Time, grouping TimeGrouping) (time.Time, error) {
	if err := grouping.validate(); err != nil {
		return time.Time{}, err
	}

	switch grouping.IntervalType {
	case IntervalSeconds:
		second := date.Second() - date.Second()%grouping.Step
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), second, 0, date.Location()), nil
	case IntervalMinutes:
		minute := date.Minute() - date.Minute()%grouping.Step
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), minute, 0, 0, date.Location()), nil
	case IntervalHours:
		hour := date.Hour() - date.Hour()%grouping.Step
		return time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location()), nil
	case IntervalDays:
		epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, date.Location())
		daysSinceEpoch := int(date.Sub(epoch) / (24 * time.Hour))
		bucketDay := (daysSinceEpoch / grouping.Step) * grouping.Step
		s := epoch.AddDate(0, 0, bucketDay)
		return time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, date.Location()), nil
	case IntervalWeeks:
		epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, date.Location())
		daysSinceEpoch := int(date.Sub(epoch) / (24 * time.Hour))
		bucketWeek := (daysSinceEpoch / 7 / grouping.Step) * grouping.Step
		s := epoch.AddDate(0, 0, bucketWeek*7)
		return time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, date.Location()), nil
	case IntervalMonths:
		month := ((int(date.Month()) - 1) / grouping.Step * grouping.Step) + 1
		return time.Date(date.Year(), time.Month(month), 1, 0, 0, 0, 0, date.Location()), nil
	case IntervalYears:
		year := ((date.Year() - 1) / grouping.Step * grouping.Step) + 1
		return time.Date(year, time.January, 1, 0, 0, 0, 0, date.Location()), nil
	default:
		return time.Time{}, fmt.Errorf("unknown interval type: %q", grouping.IntervalType)
	}
}

func addMonthsWithOverflowClamp(date time.Time, step int) time.Time {
	originalDay := date.Day()
	candidate := date.AddDate(0, step, 0)
	if candidate.Day() != originalDay {
		return lastDayOfPreviousMonth(candidate)
	}
	return candidate
}

func addYearsWithOverflowClamp(date time.Time, step int) time.Time {
	originalMonth := date.Month()
	candidate := date.AddDate(step, 0, 0)
	if candidate.Month() != originalMonth {
		return lastDayOfPreviousMonth(candidate)
	}
	return candidate
}

// lastDayOfPreviousMonth returns the last day of the month preceding date.
// The time-of-day (hour, minute, second, nanosecond) is preserved from the input
// so that AddInterval results remain on the same wall-clock time within the prior month.
func lastDayOfPreviousMonth(date time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month(),
		1,
		date.Hour(),
		date.Minute(),
		date.Second(),
		date.Nanosecond(),
		date.Location(),
	).AddDate(0, 0, -1)
}

// GenerateTimeSeries creates a series of times from start to end (inclusive) advancing
// by step intervalType units at each step.
//
// The result slice cannot be pre-allocated: month and year intervals use overflow
// clamping, so consecutive steps may land on the same date, making the final count
// unknowable without iterating.
func GenerateTimeSeries(start, end time.Time, step int, intervalType IntervalType) []time.Time {
	series := make([]time.Time, 0)
	current := start
	for !current.After(end) {
		series = append(series, current)
		current = AddInterval(current, step, intervalType)
	}
	return series
}

// SliceByTimeRange filters sequence by inclusive time bounds.
func SliceByTimeRange(sequence []SeqFactor, start, end time.Time) []SeqFactor {
	result := make([]SeqFactor, 0)
	for _, f := range sequence {
		if !f.Time.Before(start) && !f.Time.After(end) {
			result = append(result, f)
		}
	}
	return result
}

// CreateCycles creates repeated factors from start to end using a fixed interval.
//
// For IntervalMonths and IntervalYears, overflow clamping applies: a cycle starting
// on the 31st will drift after the first short month (see AddInterval). This is
// intentional JS-like behaviour.
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
