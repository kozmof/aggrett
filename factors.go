package aggrett

import (
	"fmt"
	"time"
)

// InsertFactor appends a new factor without mutating the original slice.
func InsertFactor(
	sequence []SeqFactor,
	tag string,
	timeValue time.Time,
	value float64,
	factor Factor,
	genID func() string,
) ([]SeqFactor, error) {
	if !factor.IsValid() {
		return nil, fmt.Errorf("invalid factor %q", factor)
	}
	result := make([]SeqFactor, 0, len(sequence)+1)
	result = append(result, sequence...)
	result = append(result, SeqFactor{ID: genID(), Tag: tag, Time: timeValue, Value: value, Factor: factor})
	return result, nil
}

// RemoveFactor removes factors by IDs.
func RemoveFactor(sequence []SeqFactor, ids []string) []SeqFactor {
	idSet := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	result := make([]SeqFactor, 0, len(sequence))
	for _, f := range sequence {
		if _, remove := idSet[f.ID]; !remove {
			result = append(result, f)
		}
	}
	return result
}

// SeqFactorUpdate models partial updates for UpdateFactor.
type SeqFactorUpdate struct {
	Tag    *string
	Time   *time.Time
	Value  *float64
	Factor *Factor
}

// UpdateFactor updates one factor by ID and returns a new slice.
func UpdateFactor(sequence []SeqFactor, id string, fields SeqFactorUpdate) ([]SeqFactor, error) {
	if fields.Factor != nil && !fields.Factor.IsValid() {
		return nil, fmt.Errorf("invalid factor %q", *fields.Factor)
	}
	result := make([]SeqFactor, 0, len(sequence))
	for _, f := range sequence {
		if f.ID != id {
			result = append(result, f)
			continue
		}

		updated := f
		if fields.Tag != nil {
			updated.Tag = *fields.Tag
		}
		if fields.Time != nil {
			updated.Time = *fields.Time
		}
		if fields.Value != nil {
			updated.Value = *fields.Value
		}
		if fields.Factor != nil {
			updated.Factor = *fields.Factor
		}
		result = append(result, updated)
	}
	return result, nil
}

// FindByID returns the first factor with the given ID and true, or a zero value and false.
func FindByID(sequence []SeqFactor, id string) (SeqFactor, bool) {
	for _, f := range sequence {
		if f.ID == id {
			return f, true
		}
	}
	return SeqFactor{}, false
}

// MergeSequences concatenates two sequences.
func MergeSequences(a, b []SeqFactor) []SeqFactor {
	result := make([]SeqFactor, 0, len(a)+len(b))
	result = append(result, a...)
	result = append(result, b...)
	return result
}
