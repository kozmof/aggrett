package aggrett

import "testing"

func TestBuildIndex(t *testing.T) {
	t.Run("empty sequence produces empty index", func(t *testing.T) {
		idx := BuildIndex(nil)
		_, ok := idx.Lookup(nil, "any")
		if ok {
			t.Fatalf("expected miss on empty index")
		}
	})

	t.Run("lookup hit returns correct factor", func(t *testing.T) {
		seq := []SeqFactor{
			makeSeqFactor(t, "a", "rent", "2024-01-01", 100, FactorMinus),
			makeSeqFactor(t, "b", "food", "2024-01-05", 50, FactorPlus),
			makeSeqFactor(t, "c", "salary", "2024-01-10", 5000, FactorPlus),
		}
		idx := BuildIndex(seq)

		f, ok := idx.Lookup(seq, "b")
		if !ok || f.Tag != "food" || f.Value != 50 {
			t.Fatalf("unexpected result ok=%v factor=%#v", ok, f)
		}
	})

	t.Run("lookup miss returns false", func(t *testing.T) {
		seq := []SeqFactor{makeSeqFactor(t, "a", "rent", "2024-01-01", 100, "")}
		idx := BuildIndex(seq)

		_, ok := idx.Lookup(seq, "nonexistent")
		if ok {
			t.Fatalf("expected miss for unknown id")
		}
	})

	t.Run("lookup all IDs in a sequence", func(t *testing.T) {
		seq := []SeqFactor{
			makeSeqFactor(t, "x", "rent", "2024-01-01", 100, ""),
			makeSeqFactor(t, "y", "food", "2024-01-02", 50, ""),
			makeSeqFactor(t, "z", "salary", "2024-01-03", 5000, ""),
		}
		idx := BuildIndex(seq)
		for _, want := range seq {
			got, ok := idx.Lookup(seq, want.ID)
			if !ok || got.ID != want.ID {
				t.Fatalf("expected hit for %q, got ok=%v", want.ID, ok)
			}
		}
	})

	t.Run("duplicate IDs: last occurrence wins", func(t *testing.T) {
		seq := []SeqFactor{
			makeSeqFactor(t, "dup", "rent", "2024-01-01", 100, ""),
			makeSeqFactor(t, "dup", "food", "2024-01-02", 50, ""),
		}
		idx := BuildIndex(seq)
		f, ok := idx.Lookup(seq, "dup")
		if !ok || f.Tag != "food" {
			t.Fatalf("expected last occurrence (food), got ok=%v tag=%q", ok, f.Tag)
		}
	})

	t.Run("stale index returns false after slice replaced", func(t *testing.T) {
		seq := []SeqFactor{makeSeqFactor(t, "a", "rent", "2024-01-01", 100, "")}
		idx := BuildIndex(seq)
		// Replace slice with different content at same length.
		newSeq := []SeqFactor{makeSeqFactor(t, "b", "food", "2024-01-02", 50, "")}
		_, ok := idx.Lookup(newSeq, "a")
		if ok {
			t.Fatalf("stale index should return false when slice content changed")
		}
	})
}
