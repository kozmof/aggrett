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

// String implements fmt.Stringer.
func (f Factor) String() string { return string(f) }

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

// Get returns the BreakdownEntry for the given tag and true, or a zero value and false.
// It is safe to call on a nil Breakdown.
func (b Breakdown) Get(tag string) (BreakdownEntry, bool) {
	if b == nil {
		return BreakdownEntry{}, false
	}
	e, ok := b[tag]
	return e, ok
}

// AccumCore is an accumulated value at a time point without tag breakdown metadata.
type AccumCore struct {
	IDs   []string
	Time  time.Time
	Store float64
}

// GetIDs returns the IDs of the factors contributing to this accumulation entry.
func (a AccumCore) GetIDs() []string { return a.IDs }

// GetTime returns the time of this accumulation entry.
func (a AccumCore) GetTime() time.Time { return a.Time }

// GetStore returns the accumulated store value at this time point.
func (a AccumCore) GetStore() float64 { return a.Store }

// Accum is an accumulated value at a time point with tag and breakdown metadata.
type Accum struct {
	AccumCore
	Tags      []string
	Breakdown Breakdown
}

// GetIDs returns the IDs of the factors contributing to this accumulation entry.
func (a Accum) GetIDs() []string { return a.IDs }

// GetTime returns the time of this accumulation entry.
func (a Accum) GetTime() time.Time { return a.Time }

// GetStore returns the accumulated store value at this time point.
func (a Accum) GetStore() float64 { return a.Store }

// Accumulated is implemented by both AccumCore and Accum.
type Accumulated interface {
	GetIDs() []string
	GetTime() time.Time
	GetStore() float64
}

// Compile-time interface checks.
var _ Accumulated = AccumCore{}
var _ Accumulated = Accum{}
