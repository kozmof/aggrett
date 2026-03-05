package aggrett

import "time"

// Factor determines whether a value is added to or removed from a running store.
type Factor string

const (
	FactorPlus  Factor = "plus"
	FactorMinus Factor = "minus"
)

// IsValid reports whether f is a known Factor.
func (f Factor) IsValid() bool {
	return f == FactorPlus || f == FactorMinus
}

// SeqFactor is a single tagged value change at a specific time.
type SeqFactor struct {
	ID     string
	Tag    string
	Time   time.Time
	Value  float64
	Factor Factor
}

// BreakdownEntry stores per-tag accumulation for a single time bucket.
// Delta is the net contribution of this tag within that bucket, not a cumulative total.
type BreakdownEntry struct {
	Delta float64
	IDs   []string
}

// Breakdown maps tags to per-tag stores for one accumulation entry.
type Breakdown map[string]BreakdownEntry

// AccumCore is an accumulated value at a time point without tag breakdown metadata.
type AccumCore struct {
	IDs   []string
	Time  time.Time
	Store float64
}

// Accum is an accumulated value at a time point with tag and breakdown metadata.
type Accum struct {
	IDs       []string
	Tags      []string
	Time      time.Time
	Store     float64
	Breakdown Breakdown
}
