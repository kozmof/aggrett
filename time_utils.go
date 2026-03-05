package aggrep

import "time"

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
