package aggrep

import "testing"

func TestAddInterval(t *testing.T) {
	t.Run("adds seconds", func(t *testing.T) {
		date := mustParseDate(t, "2024-01-01T10:00:00")
		result := AddInterval(date, 30, IntervalSeconds)
		mustTimeEqual(t, result, mustParseDate(t, "2024-01-01T10:00:30"))
	})

	t.Run("adds minutes", func(t *testing.T) {
		date := mustParseDate(t, "2024-01-01T10:00:00")
		result := AddInterval(date, 15, IntervalMinutes)
		mustTimeEqual(t, result, mustParseDate(t, "2024-01-01T10:15:00"))
	})

	t.Run("adds hours", func(t *testing.T) {
		date := mustParseDate(t, "2024-01-01T10:00:00")
		result := AddInterval(date, 3, IntervalHours)
		mustTimeEqual(t, result, mustParseDate(t, "2024-01-01T13:00:00"))
	})

	t.Run("adds days", func(t *testing.T) {
		date := mustParseDate(t, "2024-01-01")
		result := AddInterval(date, 5, IntervalDays)
		mustTimeEqual(t, result, mustParseDate(t, "2024-01-06"))
	})

	t.Run("adds weeks", func(t *testing.T) {
		date := mustParseDate(t, "2024-01-01")
		result := AddInterval(date, 2, IntervalWeeks)
		mustTimeEqual(t, result, mustParseDate(t, "2024-01-15"))
	})

	t.Run("adds months", func(t *testing.T) {
		date := mustParseDate(t, "2024-01-15")
		result := AddInterval(date, 1, IntervalMonths)
		mustTimeEqual(t, result, mustParseDate(t, "2024-02-15"))
	})

	t.Run("handles month overflow", func(t *testing.T) {
		date := mustParseDate(t, "2024-01-31")
		result := AddInterval(date, 1, IntervalMonths)
		mustTimeEqual(t, result, mustParseDate(t, "2024-02-29"))
	})

	t.Run("adds years", func(t *testing.T) {
		date := mustParseDate(t, "2024-06-01")
		result := AddInterval(date, 1, IntervalYears)
		mustTimeEqual(t, result, mustParseDate(t, "2025-06-01"))
	})

	t.Run("handles leap year", func(t *testing.T) {
		date := mustParseDate(t, "2024-02-29")
		result := AddInterval(date, 1, IntervalYears)
		mustTimeEqual(t, result, mustParseDate(t, "2025-02-28"))
	})

	t.Run("does not mutate original date", func(t *testing.T) {
		date := mustParseDate(t, "2024-01-01")
		_ = AddInterval(date, 5, IntervalDays)
		mustTimeEqual(t, date, mustParseDate(t, "2024-01-01"))
	})
}

func TestGenerateTimeSeries(t *testing.T) {
	t.Run("generates daily series", func(t *testing.T) {
		start := mustParseDate(t, "2024-01-01")
		end := mustParseDate(t, "2024-01-03")
		result := GenerateTimeSeries(start, end, 1, IntervalDays)

		if len(result) != 3 {
			t.Fatalf("got %d items want 3", len(result))
		}
		mustTimeEqual(t, result[0], mustParseDate(t, "2024-01-01"))
		mustTimeEqual(t, result[1], mustParseDate(t, "2024-01-02"))
		mustTimeEqual(t, result[2], mustParseDate(t, "2024-01-03"))
	})

	t.Run("generates with step greater than 1", func(t *testing.T) {
		start := mustParseDate(t, "2024-01-01")
		end := mustParseDate(t, "2024-01-10")
		result := GenerateTimeSeries(start, end, 3, IntervalDays)

		if len(result) != 4 {
			t.Fatalf("got %d items want 4", len(result))
		}
		mustTimeEqual(t, result[0], mustParseDate(t, "2024-01-01"))
		mustTimeEqual(t, result[1], mustParseDate(t, "2024-01-04"))
		mustTimeEqual(t, result[2], mustParseDate(t, "2024-01-07"))
		mustTimeEqual(t, result[3], mustParseDate(t, "2024-01-10"))
	})

	t.Run("returns single item when start equals end", func(t *testing.T) {
		date := mustParseDate(t, "2024-01-01")
		result := GenerateTimeSeries(date, date, 1, IntervalDays)

		if len(result) != 1 {
			t.Fatalf("got %d items want 1", len(result))
		}
		mustTimeEqual(t, result[0], date)
	})

	t.Run("returns empty when start greater than end", func(t *testing.T) {
		result := GenerateTimeSeries(mustParseDate(t, "2024-02-01"), mustParseDate(t, "2024-01-01"), 1, IntervalDays)
		if len(result) != 0 {
			t.Fatalf("got %d items want 0", len(result))
		}
	})
}

func TestSliceByTimeRange(t *testing.T) {
	seq := []SeqFactor{
		{ID: "a", Tag: "rent", Time: mustParseDate(t, "2024-01-01"), Value: 100, Factor: FactorPlus},
		{ID: "b", Tag: "food", Time: mustParseDate(t, "2024-02-01"), Value: 50, Factor: FactorPlus},
		{ID: "c", Tag: "rent", Time: mustParseDate(t, "2024-03-01"), Value: 100, Factor: FactorPlus},
		{ID: "d", Tag: "food", Time: mustParseDate(t, "2024-04-01"), Value: 50, Factor: FactorPlus},
	}

	t.Run("returns factors within the time range inclusive", func(t *testing.T) {
		result := SliceByTimeRange(seq, mustParseDate(t, "2024-02-01"), mustParseDate(t, "2024-03-01"))
		if len(result) != 2 || result[0].ID != "b" || result[1].ID != "c" {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("includes factors exactly at boundaries", func(t *testing.T) {
		result := SliceByTimeRange(seq, mustParseDate(t, "2024-01-01"), mustParseDate(t, "2024-04-01"))
		if len(result) != 4 {
			t.Fatalf("got %d items want 4", len(result))
		}
	})

	t.Run("returns empty when no factors in range", func(t *testing.T) {
		result := SliceByTimeRange(seq, mustParseDate(t, "2025-01-01"), mustParseDate(t, "2025-12-31"))
		if len(result) != 0 {
			t.Fatalf("got %d items want 0", len(result))
		}
	})

	t.Run("does not mutate original array", func(t *testing.T) {
		_ = SliceByTimeRange(seq, mustParseDate(t, "2024-02-01"), mustParseDate(t, "2024-03-01"))
		if len(seq) != 4 {
			t.Fatalf("input sequence mutated")
		}
	})

	t.Run("works on empty sequence", func(t *testing.T) {
		result := SliceByTimeRange(nil, mustParseDate(t, "2024-01-01"), mustParseDate(t, "2024-12-31"))
		if len(result) != 0 {
			t.Fatalf("got %d items want 0", len(result))
		}
	})
}
