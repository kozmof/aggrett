package aggrett

import (
	"reflect"
	"testing"
)

func makeTagFactor(t *testing.T, id, tag, timeValue string, value float64, factor Factor) SeqFactor {
	t.Helper()
	if factor == "" {
		factor = FactorPlus
	}
	return SeqFactor{ID: id, Tag: tag, Time: mustParseDate(t, timeValue), Value: value, Factor: factor}
}

func TestFilterByTag(t *testing.T) {
	seq := []SeqFactor{
		makeTagFactor(t, "1", "rent", "2024-01-01", 1000, FactorMinus),
		makeTagFactor(t, "2", "food", "2024-01-05", 200, FactorMinus),
		makeTagFactor(t, "3", "salary", "2024-01-10", 5000, FactorPlus),
		makeTagFactor(t, "4", "food", "2024-01-15", 150, FactorMinus),
		makeTagFactor(t, "5", "utilities", "2024-01-20", 100, FactorMinus),
		makeTagFactor(t, "6", "rent", "2024-02-01", 1000, FactorMinus),
	}

	t.Run("filters to matching tags", func(t *testing.T) {
		result := FilterByTag(seq, []string{"food"})
		if len(result) != 2 {
			t.Fatalf("got %d items want 2", len(result))
		}
		for _, f := range result {
			if f.Tag != "food" {
				t.Fatalf("unexpected tag %q", f.Tag)
			}
		}
	})

	t.Run("filters to multiple tags", func(t *testing.T) {
		result := FilterByTag(seq, []string{"rent", "utilities"})
		if len(result) != 3 || result[0].ID != "1" || result[1].ID != "5" || result[2].ID != "6" {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("returns empty array when no tags match", func(t *testing.T) {
		if len(FilterByTag(seq, []string{"nonexistent"})) != 0 {
			t.Fatalf("expected empty result")
		}
	})

	t.Run("returns empty array for empty sequence", func(t *testing.T) {
		if len(FilterByTag(nil, []string{"food"})) != 0 {
			t.Fatalf("expected empty result")
		}
	})

	t.Run("does not mutate the original array", func(t *testing.T) {
		_ = FilterByTag(seq, []string{"food"})
		if len(seq) != 6 {
			t.Fatalf("input sequence mutated")
		}
	})
}

func TestExtractTags(t *testing.T) {
	seq := []SeqFactor{
		makeTagFactor(t, "1", "rent", "2024-01-01", 1000, FactorMinus),
		makeTagFactor(t, "2", "food", "2024-01-05", 200, FactorMinus),
		makeTagFactor(t, "3", "salary", "2024-01-10", 5000, FactorPlus),
		makeTagFactor(t, "4", "food", "2024-01-15", 150, FactorMinus),
		makeTagFactor(t, "5", "utilities", "2024-01-20", 100, FactorMinus),
		makeTagFactor(t, "6", "rent", "2024-02-01", 1000, FactorMinus),
	}

	t.Run("returns unique tags", func(t *testing.T) {
		tags := ExtractTags(seq)
		if len(tags) != 4 {
			t.Fatalf("got %d tags want 4", len(tags))
		}
		if !reflect.DeepEqual(tags, []string{"rent", "food", "salary", "utilities"}) {
			t.Fatalf("unexpected tags %#v", tags)
		}
	})

	t.Run("returns empty array for empty sequence", func(t *testing.T) {
		if len(ExtractTags(nil)) != 0 {
			t.Fatalf("expected empty tags")
		}
	})

	t.Run("returns single tag when all items share it", func(t *testing.T) {
		uniform := []SeqFactor{
			makeTagFactor(t, "1", "rent", "2024-01-01", 100, ""),
			makeTagFactor(t, "2", "rent", "2024-02-01", 100, ""),
		}
		tags := ExtractTags(uniform)
		if !reflect.DeepEqual(tags, []string{"rent"}) {
			t.Fatalf("unexpected tags %#v", tags)
		}
	})
}

func TestGroupByTag(t *testing.T) {
	seq := []SeqFactor{
		makeTagFactor(t, "1", "rent", "2024-01-01", 1000, FactorMinus),
		makeTagFactor(t, "2", "food", "2024-01-05", 200, FactorMinus),
		makeTagFactor(t, "3", "salary", "2024-01-10", 5000, FactorPlus),
		makeTagFactor(t, "4", "food", "2024-01-15", 150, FactorMinus),
		makeTagFactor(t, "5", "utilities", "2024-01-20", 100, FactorMinus),
		makeTagFactor(t, "6", "rent", "2024-02-01", 1000, FactorMinus),
	}

	t.Run("groups factors by tag", func(t *testing.T) {
		groups := GroupByTag(seq)
		if len(groups) != 4 {
			t.Fatalf("got %d groups want 4", len(groups))
		}
		if len(groups["rent"]) != 2 || len(groups["food"]) != 2 || len(groups["salary"]) != 1 || len(groups["utilities"]) != 1 {
			t.Fatalf("unexpected groups %#v", groups)
		}
	})

	t.Run("preserves factor data in groups", func(t *testing.T) {
		groups := GroupByTag(seq)
		if groups["salary"][0].ID != "3" || groups["salary"][0].Value != 5000 {
			t.Fatalf("unexpected salary group %#v", groups["salary"])
		}
	})

	t.Run("returns empty object for empty sequence", func(t *testing.T) {
		groups := GroupByTag(nil)
		if len(groups) != 0 {
			t.Fatalf("expected empty groups")
		}
	})
}

func TestExcludeByTag(t *testing.T) {
	seq := []SeqFactor{
		makeTagFactor(t, "1", "rent", "2024-01-01", 1000, FactorMinus),
		makeTagFactor(t, "2", "food", "2024-01-05", 200, FactorMinus),
		makeTagFactor(t, "3", "salary", "2024-01-10", 5000, FactorPlus),
		makeTagFactor(t, "4", "food", "2024-01-15", 150, FactorMinus),
		makeTagFactor(t, "5", "utilities", "2024-01-20", 100, FactorMinus),
		makeTagFactor(t, "6", "rent", "2024-02-01", 1000, FactorMinus),
	}

	t.Run("excludes matching tags", func(t *testing.T) {
		result := ExcludeByTag(seq, []string{"food"})
		if len(result) != 4 {
			t.Fatalf("got %d items want 4", len(result))
		}
		for _, f := range result {
			if f.Tag == "food" {
				t.Fatalf("food tag should have been excluded")
			}
		}
	})

	t.Run("excludes multiple tags", func(t *testing.T) {
		result := ExcludeByTag(seq, []string{"food", "utilities"})
		if len(result) != 3 || result[0].ID != "1" || result[1].ID != "3" || result[2].ID != "6" {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("returns all when no tags match", func(t *testing.T) {
		result := ExcludeByTag(seq, []string{"nonexistent"})
		if len(result) != 6 {
			t.Fatalf("got %d items want 6", len(result))
		}
	})

	t.Run("returns empty array for empty sequence", func(t *testing.T) {
		if len(ExcludeByTag(nil, []string{"food"})) != 0 {
			t.Fatalf("expected empty result")
		}
	})

	t.Run("does not mutate the original array", func(t *testing.T) {
		_ = ExcludeByTag(seq, []string{"food"})
		if len(seq) != 6 {
			t.Fatalf("input sequence mutated")
		}
	})
}

func TestRemoveByTag(t *testing.T) {
	seq := []SeqFactor{
		makeTagFactor(t, "1", "rent", "2024-01-01", 1000, FactorMinus),
		makeTagFactor(t, "2", "food", "2024-01-05", 200, FactorMinus),
		makeTagFactor(t, "3", "salary", "2024-01-10", 5000, FactorPlus),
	}

	excluded := ExcludeByTag(seq, []string{"food", "rent"})
	removed := RemoveByTag(seq, []string{"food", "rent"})
	if !reflect.DeepEqual(removed, excluded) {
		t.Fatalf("removeByTag mismatch: removed=%#v excluded=%#v", removed, excluded)
	}
}

func TestRenameTag(t *testing.T) {
	seq := []SeqFactor{
		makeTagFactor(t, "1", "rent", "2024-01-01", 1000, FactorMinus),
		makeTagFactor(t, "2", "food", "2024-01-05", 200, FactorMinus),
		makeTagFactor(t, "3", "salary", "2024-01-10", 5000, FactorPlus),
		makeTagFactor(t, "4", "food", "2024-01-15", 150, FactorMinus),
		makeTagFactor(t, "5", "utilities", "2024-01-20", 100, FactorMinus),
		makeTagFactor(t, "6", "rent", "2024-02-01", 1000, FactorMinus),
	}

	t.Run("renames all matching factors", func(t *testing.T) {
		result := RenameTag(seq, "food", "dining")
		diningCount := 0
		foodCount := 0
		for _, f := range result {
			if f.Tag == "dining" {
				diningCount++
			}
			if f.Tag == "food" {
				foodCount++
			}
		}
		if diningCount != 2 || foodCount != 0 {
			t.Fatalf("unexpected tag counts dining=%d food=%d", diningCount, foodCount)
		}
	})

	t.Run("preserves other fields", func(t *testing.T) {
		result := RenameTag(seq, "food", "dining")
		var renamed SeqFactor
		for _, f := range result {
			if f.ID == "2" {
				renamed = f
				break
			}
		}
		if renamed.Tag != "dining" || renamed.Value != 200 || renamed.Factor != FactorMinus {
			t.Fatalf("unexpected renamed factor %#v", renamed)
		}
	})

	t.Run("leaves non-matching factors unchanged", func(t *testing.T) {
		result := RenameTag(seq, "food", "dining")
		for _, f := range result {
			if f.ID == "1" && f.Tag != "rent" {
				t.Fatalf("non-matching factor changed %#v", f)
			}
		}
	})

	t.Run("returns unchanged copy when tag not found", func(t *testing.T) {
		result := RenameTag(seq, "nonexistent", "new")
		if len(result) != 6 {
			t.Fatalf("got %d items want 6", len(result))
		}
		if !reflect.DeepEqual(result, seq) {
			t.Fatalf("expected unchanged copy")
		}
		if &result[0] == &seq[0] {
			t.Fatalf("expected copy, got same backing elements")
		}
	})

	t.Run("does not mutate the original array", func(t *testing.T) {
		_ = RenameTag(seq, "food", "dining")
		if seq[1].Tag != "food" {
			t.Fatalf("input sequence mutated")
		}
	})
}

func TestAccumulateByTag(t *testing.T) {
	seq := []SeqFactor{
		makeTagFactor(t, "1", "rent", "2024-01-01", 1000, FactorMinus),
		makeTagFactor(t, "2", "food", "2024-01-05", 200, FactorMinus),
		makeTagFactor(t, "3", "salary", "2024-01-10", 5000, FactorPlus),
		makeTagFactor(t, "4", "food", "2024-01-15", 150, FactorMinus),
		makeTagFactor(t, "5", "utilities", "2024-01-20", 100, FactorMinus),
		makeTagFactor(t, "6", "rent", "2024-02-01", 1000, FactorMinus),
	}

	t.Run("accumulates only the specified tag", func(t *testing.T) {
		result := AccumulateByTag(seq, 0, "food")
		if len(result) != 2 || result[0].Store != -200 || result[1].Store != -350 {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("returns empty array when tag not found", func(t *testing.T) {
		result := AccumulateByTag(seq, 0, "nonexistent")
		if len(result) != 0 {
			t.Fatalf("expected empty result")
		}
	})

	t.Run("applies base value", func(t *testing.T) {
		result := AccumulateByTag(seq, 1000, "rent")
		if len(result) != 2 || result[0].Store != 0 || result[1].Store != -1000 {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("returns empty array for empty sequence", func(t *testing.T) {
		result := AccumulateByTag(nil, 0, "food")
		if len(result) != 0 {
			t.Fatalf("expected empty result")
		}
	})
}
