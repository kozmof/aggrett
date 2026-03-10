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

// timeGrouping describes interval-based bucketing for accumulation and aggregation.
// Result times are the start of each bucket.
type timeGrouping struct {
	Step         int
	IntervalType IntervalType
}

func (g timeGrouping) validate() {
	if g.Step <= 0 {
		panic(fmt.Sprintf("grouping step must be positive: %d", g.Step))
	}
	if !g.IntervalType.IsValid() {
		panic(fmt.Sprintf("unknown interval type: %q", g.IntervalType))
	}
}

// IsValid reports whether t is a known IntervalType.
func (t IntervalType) IsValid() bool {
	switch t {
	case IntervalSeconds, IntervalMinutes, IntervalHours,
		IntervalDays, IntervalWeeks, IntervalMonths, IntervalYears:
		return true
	}
	return false
}

// AddInterval adds an interval to a date with JS-like month/year overflow behavior.
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

func bucketStart(date time.Time, grouping timeGrouping) time.Time {
	grouping.validate()

	switch grouping.IntervalType {
	case IntervalSeconds:
		second := date.Second() - date.Second()%grouping.Step
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), second, 0, date.Location())
	case IntervalMinutes:
		minute := date.Minute() - date.Minute()%grouping.Step
		return time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), minute, 0, 0, date.Location())
	case IntervalHours:
		hour := date.Hour() - date.Hour()%grouping.Step
		return time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
	case IntervalDays:
		day := ((date.Day() - 1) / grouping.Step * grouping.Step) + 1
		return time.Date(date.Year(), date.Month(), day, 0, 0, 0, 0, date.Location())
	case IntervalWeeks:
		offset := ((date.YearDay() - 1) / (grouping.Step * 7)) * grouping.Step * 7
		start := time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, date.Location())
		return start.AddDate(0, 0, offset)
	case IntervalMonths:
		month := ((int(date.Month()) - 1) / grouping.Step * grouping.Step) + 1
		return time.Date(date.Year(), time.Month(month), 1, 0, 0, 0, 0, date.Location())
	case IntervalYears:
		year := ((date.Year() - 1) / grouping.Step * grouping.Step) + 1
		return time.Date(year, time.January, 1, 0, 0, 0, 0, date.Location())
	default:
		panic(fmt.Sprintf("unknown interval type: %q", grouping.IntervalType))
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

// GenerateTimeSeries creates a series from start to end (inclusive).
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
