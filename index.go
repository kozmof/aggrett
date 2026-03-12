package aggrett

// Index is an optional companion to a []SeqFactor slice that enables O(1) lookup
// by ID. It must be rebuilt whenever the underlying slice changes via BuildIndex.
//
// The Index stores positions (not pointers) so it remains valid when the slice is
// copied, as long as the slice contents have not changed.
type Index struct {
	m map[string]int // ID -> position in the slice
}

// BuildIndex constructs an Index from sequence. If duplicate IDs exist, the last
// occurrence wins.
func BuildIndex(sequence []SeqFactor) Index {
	m := make(map[string]int, len(sequence))
	for i, f := range sequence {
		m[f.ID] = i
	}
	return Index{m: m}
}

// Lookup returns the SeqFactor at the indexed position and true if the ID is found,
// or a zero value and false otherwise. sequence must be the same slice that was
// passed to BuildIndex.
func (idx Index) Lookup(sequence []SeqFactor, id string) (SeqFactor, bool) {
	i, ok := idx.m[id]
	if !ok || i >= len(sequence) {
		return SeqFactor{}, false
	}
	f := sequence[i]
	if f.ID != id {
		// Position is stale (slice was mutated after index was built).
		return SeqFactor{}, false
	}
	return f, true
}
