package aggrett

import (
	"fmt"
	"strings"
	"testing"
)

func TestInsertFactor(t *testing.T) {
	mockGenID := func() func() string {
		counter := 0
		return func() string {
			counter++
			return fmt.Sprintf("fac-%d", counter)
		}
	}()

	t.Run("appends a new factor to the sequence", func(t *testing.T) {
		seq := []SeqFactor{makeSeqFactor(t, "a", "rent", "2024-01-01", 100, "")}
		result, err := InsertFactor(seq, "food", mustParseDate(t, "2024-02-01"), 50, FactorMinus, mockGenID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 2 {
			t.Fatalf("got %d items want 2", len(result))
		}
		if result[1].Tag != "food" || result[1].Value != 50 || result[1].Factor != FactorMinus {
			t.Fatalf("unexpected inserted factor %#v", result[1])
		}
		if !strings.HasPrefix(result[1].ID, "fac-") {
			t.Fatalf("unexpected id %q", result[1].ID)
		}
	})

	t.Run("does not mutate original array", func(t *testing.T) {
		seq := []SeqFactor{makeSeqFactor(t, "a", "rent", "2024-01-01", 100, "")}
		_, _ = InsertFactor(seq, "food", mustParseDate(t, "2024-02-01"), 50, FactorPlus, mockGenID)
		if len(seq) != 1 {
			t.Fatalf("input sequence mutated")
		}
	})

	t.Run("works on an empty sequence", func(t *testing.T) {
		result, err := InsertFactor(nil, "salary", mustParseDate(t, "2024-03-01"), 3000, FactorPlus, mockGenID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 || result[0].Tag != "salary" {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("returns error for invalid factor", func(t *testing.T) {
		_, err := InsertFactor(nil, "salary", mustParseDate(t, "2024-03-01"), 3000, Factor("bad"), mockGenID)
		if err == nil {
			t.Fatalf("expected error for invalid factor")
		}
	})

	t.Run("returns error when genID produces a duplicate ID", func(t *testing.T) {
		seq := []SeqFactor{makeSeqFactor(t, "dup-id", "rent", "2024-01-01", 100, "")}
		alwaysDup := func() string { return "dup-id" }
		_, err := InsertFactor(seq, "food", mustParseDate(t, "2024-02-01"), 50, FactorPlus, alwaysDup)
		if err == nil {
			t.Fatalf("expected error for duplicate ID")
		}
	})
}

func TestRemoveFactor(t *testing.T) {
	seq := []SeqFactor{
		makeSeqFactor(t, "a", "rent", "2024-01-01", 100, ""),
		makeSeqFactor(t, "b", "food", "2024-01-15", 50, ""),
		makeSeqFactor(t, "c", "rent", "2024-02-01", 100, ""),
	}

	t.Run("removes a single factor by ID", func(t *testing.T) {
		result := RemoveFactor(seq, []string{"b"})
		if len(result) != 2 || result[0].ID != "a" || result[1].ID != "c" {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("removes multiple factors by ID", func(t *testing.T) {
		result := RemoveFactor(seq, []string{"a", "c"})
		if len(result) != 1 || result[0].ID != "b" {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("returns unchanged copy when ID not found", func(t *testing.T) {
		result := RemoveFactor(seq, []string{"nonexistent"})
		if len(result) != 3 {
			t.Fatalf("got %d items want 3", len(result))
		}
		if &result[0] == &seq[0] {
			t.Fatalf("expected a copy, got same backing elements")
		}
	})

	t.Run("does not mutate original array", func(t *testing.T) {
		_ = RemoveFactor(seq, []string{"a"})
		if len(seq) != 3 {
			t.Fatalf("input sequence mutated")
		}
	})

	t.Run("works on an empty sequence", func(t *testing.T) {
		result := RemoveFactor(nil, []string{"a"})
		if len(result) != 0 {
			t.Fatalf("got %d items want 0", len(result))
		}
	})
}

func TestUpdateFactor(t *testing.T) {
	seq := []SeqFactor{
		makeSeqFactor(t, "a", "rent", "2024-01-01", 100, FactorMinus),
		makeSeqFactor(t, "b", "food", "2024-01-15", 50, FactorMinus),
	}

	t.Run("updates a single field", func(t *testing.T) {
		result, err := UpdateFactor(seq, "a", SeqFactorUpdate{Value: floatPtr(200)})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result[0].Value != 200 || result[0].Tag != "rent" || result[0].Factor != FactorMinus {
			t.Fatalf("unexpected result %#v", result[0])
		}
	})

	t.Run("updates multiple fields", func(t *testing.T) {
		result, err := UpdateFactor(seq, "b", SeqFactorUpdate{Tag: strPtr("dining"), Factor: factorPtr(FactorPlus), Value: floatPtr(75)})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result[1].Tag != "dining" || result[1].Factor != FactorPlus || result[1].Value != 75 {
			t.Fatalf("unexpected result %#v", result[1])
		}
	})

	t.Run("returns unchanged copy when ID not found", func(t *testing.T) {
		result, err := UpdateFactor(seq, "nonexistent", SeqFactorUpdate{Value: floatPtr(999)})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 || result[0].Value != 100 || result[1].Value != 50 {
			t.Fatalf("unexpected result %#v", result)
		}
		if &result[0] == &seq[0] {
			t.Fatalf("expected a copy, got same backing elements")
		}
	})

	t.Run("does not mutate the original array", func(t *testing.T) {
		_, _ = UpdateFactor(seq, "a", SeqFactorUpdate{Value: floatPtr(999)})
		if seq[0].Value != 100 {
			t.Fatalf("input sequence mutated")
		}
	})

	t.Run("returns error for invalid factor", func(t *testing.T) {
		bad := Factor("bad")
		_, err := UpdateFactor(seq, "a", SeqFactorUpdate{Factor: &bad})
		if err == nil {
			t.Fatalf("expected error for invalid factor")
		}
	})
}

func TestMergeSequences(t *testing.T) {
	t.Run("concatenates two sequences", func(t *testing.T) {
		a := []SeqFactor{makeSeqFactor(t, "a", "rent", "2024-01-01", 100, "")}
		b := []SeqFactor{makeSeqFactor(t, "b", "food", "2024-02-01", 50, "")}
		result := MergeSequences(a, b)

		if len(result) != 2 || result[0].ID != "a" || result[1].ID != "b" {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("merging with empty array returns copy of the other", func(t *testing.T) {
		a := []SeqFactor{makeSeqFactor(t, "a", "rent", "2024-01-01", 100, "")}
		if len(MergeSequences(a, nil)) != 1 || len(MergeSequences(nil, a)) != 1 {
			t.Fatalf("unexpected merge length")
		}
	})

	t.Run("does not mutate the original arrays", func(t *testing.T) {
		a := []SeqFactor{makeSeqFactor(t, "a", "rent", "2024-01-01", 100, "")}
		b := []SeqFactor{makeSeqFactor(t, "b", "food", "2024-02-01", 50, "")}
		_ = MergeSequences(a, b)
		if len(a) != 1 || len(b) != 1 {
			t.Fatalf("input arrays mutated")
		}
	})
}
