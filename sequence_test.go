package aggrett

import "testing"

func TestAccumulate(t *testing.T) {
	t.Run("adds value with plus factor", func(t *testing.T) {
		got, err := Accumulate(FactorPlus, 10, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 15 {
			t.Fatalf("got %v want 15", got)
		}
	})

	t.Run("subtracts value with minus factor", func(t *testing.T) {
		got, err := Accumulate(FactorMinus, 10, 5)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 5 {
			t.Fatalf("got %v want 5", got)
		}
	})

	t.Run("works with zero base value", func(t *testing.T) {
		got, err := Accumulate(FactorPlus, 0, 100)
		if err != nil || got != 100 {
			t.Fatalf("got %v err %v want 100 nil", got, err)
		}
		got, err = Accumulate(FactorMinus, 0, 100)
		if err != nil || got != -100 {
			t.Fatalf("got %v err %v want -100 nil", got, err)
		}
	})

	t.Run("works with zero value", func(t *testing.T) {
		got, err := Accumulate(FactorPlus, 50, 0)
		if err != nil || got != 50 {
			t.Fatalf("got %v err %v want 50 nil", got, err)
		}
		got, err = Accumulate(FactorMinus, 50, 0)
		if err != nil || got != 50 {
			t.Fatalf("got %v err %v want 50 nil", got, err)
		}
	})

	t.Run("handles negative results", func(t *testing.T) {
		got, err := Accumulate(FactorMinus, 3, 10)
		if err != nil || got != -7 {
			t.Fatalf("got %v err %v want -7 nil", got, err)
		}
	})

	t.Run("returns error for unknown factor", func(t *testing.T) {
		_, err := Accumulate(Factor("bad"), 10, 5)
		if err == nil {
			t.Fatalf("expected error for unknown factor")
		}
	})
}

func TestAccumulateSequence(t *testing.T) {
	t.Run("returns empty array for empty sequence", func(t *testing.T) {
		result, err := AccumulateSequence(nil, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Fatalf("got %d items want 0", len(result))
		}
	})

	t.Run("accumulates a single factor", func(t *testing.T) {
		sequence := []SeqFactor{{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-01"), Factor: FactorPlus, Value: 10}}
		result, err := AccumulateSequence(sequence, 100)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 1 {
			t.Fatalf("got %d items want 1", len(result))
		}
		if result[0].Store != 110 {
			t.Fatalf("got store %v want 110", result[0].Store)
		}
		if len(result[0].IDs) != 1 || result[0].IDs[0] != "1" {
			t.Fatalf("got ids %#v want [1]", result[0].IDs)
		}
	})

	t.Run("groups factors at the same timestamp", func(t *testing.T) {
		timeValue := mustParseDate(t, "2024-01-01")
		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: timeValue, Factor: FactorPlus, Value: 10},
			{ID: "2", Tag: "b", Time: timeValue, Factor: FactorMinus, Value: 3},
		}
		result, err := AccumulateSequence(sequence, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 1 {
			t.Fatalf("got %d items want 1", len(result))
		}
		if result[0].Store != 7 {
			t.Fatalf("got store %v want 7", result[0].Store)
		}
		if len(result[0].IDs) != 2 || result[0].IDs[0] != "1" || result[0].IDs[1] != "2" {
			t.Fatalf("got ids %#v want [1 2]", result[0].IDs)
		}
	})

	t.Run("carries running total across time points", func(t *testing.T) {
		yesterday := mustParseDate(t, "2024-01-01")
		today := mustParseDate(t, "2024-01-02")
		tomorrow := mustParseDate(t, "2024-01-03")

		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: yesterday, Factor: FactorPlus, Value: 10},
			{ID: "2", Tag: "a", Time: today, Factor: FactorPlus, Value: 5},
			{ID: "3", Tag: "b", Time: tomorrow, Factor: FactorMinus, Value: 3},
		}
		result, err := AccumulateSequence(sequence, 100)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 3 {
			t.Fatalf("got %d items want 3", len(result))
		}
		if result[0].Store != 110 || result[1].Store != 115 || result[2].Store != 112 {
			t.Fatalf("unexpected stores %#v", []float64{result[0].Store, result[1].Store, result[2].Store})
		}
	})

	t.Run("sorts by time regardless of input order", func(t *testing.T) {
		earlier := mustParseDate(t, "2024-01-01")
		later := mustParseDate(t, "2024-01-02")

		sequence := []SeqFactor{
			{ID: "2", Tag: "a", Time: later, Factor: FactorPlus, Value: 5},
			{ID: "1", Tag: "a", Time: earlier, Factor: FactorPlus, Value: 10},
		}
		result, err := AccumulateSequence(sequence, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 2 {
			t.Fatalf("got %d items want 2", len(result))
		}
		if result[0].IDs[0] != "1" || result[0].Store != 10 || result[1].IDs[0] != "2" || result[1].Store != 15 {
			t.Fatalf("unexpected result %#v", result)
		}
	})

	t.Run("does not mutate input sequence", func(t *testing.T) {
		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-02"), Factor: FactorPlus, Value: 10},
			{ID: "2", Tag: "a", Time: mustParseDate(t, "2024-01-01"), Factor: FactorPlus, Value: 5},
		}
		original := append([]SeqFactor{}, sequence...)
		_, _ = AccumulateSequence(sequence, 0)

		if sequence[0].ID != original[0].ID || sequence[1].ID != original[1].ID {
			t.Fatalf("input sequence was mutated")
		}
	})
}

func TestAccumulateSequenceByInterval(t *testing.T) {
	t.Run("groups factors into 15 second buckets", func(t *testing.T) {
		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-01T10:00:14"), Factor: FactorPlus, Value: 10},
			{ID: "2", Tag: "a", Time: mustParseDate(t, "2024-01-01T10:00:01"), Factor: FactorPlus, Value: 5},
			{ID: "3", Tag: "a", Time: mustParseDate(t, "2024-01-01T10:00:15"), Factor: FactorMinus, Value: 2},
			{ID: "4", Tag: "a", Time: mustParseDate(t, "2024-01-01T10:00:29"), Factor: FactorPlus, Value: 3},
		}

		result, err := AccumulateSequenceByInterval(sequence, 0, 15, IntervalSeconds)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 2 {
			t.Fatalf("got %d items want 2", len(result))
		}
		mustTimeEqual(t, result[0].Time, mustParseDate(t, "2024-01-01T10:00:00"))
		mustTimeEqual(t, result[1].Time, mustParseDate(t, "2024-01-01T10:00:15"))
		if result[0].Store != 15 || result[1].Store != 16 {
			t.Fatalf("unexpected stores %#v", []float64{result[0].Store, result[1].Store})
		}
		if len(result[0].IDs) != 2 || result[0].IDs[0] != "2" || result[0].IDs[1] != "1" {
			t.Fatalf("unexpected first bucket ids %#v", result[0].IDs)
		}
	})

	t.Run("groups factors into 8 minute buckets", func(t *testing.T) {
		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-01T10:00:00"), Factor: FactorPlus, Value: 10},
			{ID: "2", Tag: "a", Time: mustParseDate(t, "2024-01-01T10:07:59"), Factor: FactorPlus, Value: 5},
			{ID: "3", Tag: "a", Time: mustParseDate(t, "2024-01-01T10:08:00"), Factor: FactorMinus, Value: 3},
			{ID: "4", Tag: "a", Time: mustParseDate(t, "2024-01-01T10:16:00"), Factor: FactorPlus, Value: 2},
		}

		result, err := AccumulateSequenceByInterval(sequence, 0, 8, IntervalMinutes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 3 {
			t.Fatalf("got %d items want 3", len(result))
		}
		mustTimeEqual(t, result[0].Time, mustParseDate(t, "2024-01-01T10:00:00"))
		mustTimeEqual(t, result[1].Time, mustParseDate(t, "2024-01-01T10:08:00"))
		mustTimeEqual(t, result[2].Time, mustParseDate(t, "2024-01-01T10:16:00"))
		if result[0].Store != 15 || result[1].Store != 12 || result[2].Store != 14 {
			t.Fatalf("unexpected stores %#v", []float64{result[0].Store, result[1].Store, result[2].Store})
		}
	})

	t.Run("groups factors into daily buckets", func(t *testing.T) {
		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-01T01:00:00"), Factor: FactorPlus, Value: 10},
			{ID: "2", Tag: "a", Time: mustParseDate(t, "2024-01-01T23:59:59"), Factor: FactorMinus, Value: 2},
			{ID: "3", Tag: "a", Time: mustParseDate(t, "2024-01-02T00:00:00"), Factor: FactorPlus, Value: 5},
		}

		result, err := AccumulateSequenceByInterval(sequence, 0, 1, IntervalDays)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result) != 2 {
			t.Fatalf("got %d items want 2", len(result))
		}
		mustTimeEqual(t, result[0].Time, mustParseDate(t, "2024-01-01"))
		mustTimeEqual(t, result[1].Time, mustParseDate(t, "2024-01-02"))
		if result[0].Store != 8 || result[1].Store != 13 {
			t.Fatalf("unexpected stores %#v", []float64{result[0].Store, result[1].Store})
		}
	})

	t.Run("returns error for non-positive step", func(t *testing.T) {
		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-01"), Factor: FactorPlus, Value: 10},
		}
		_, err := AccumulateSequenceByInterval(sequence, 0, 0, IntervalDays)
		if err == nil {
			t.Fatalf("expected error for step=0")
		}
	})

	t.Run("returns error for invalid interval type", func(t *testing.T) {
		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-01"), Factor: FactorPlus, Value: 10},
		}
		_, err := AccumulateSequenceByInterval(sequence, 0, 1, IntervalType("bad"))
		if err == nil {
			t.Fatalf("expected error for invalid interval type")
		}
	})

	t.Run("step=3 days: Jan-31 and Feb-1 land in same epoch-aligned bucket", func(t *testing.T) {
		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-31"), Factor: FactorPlus, Value: 10},
			{ID: "2", Tag: "a", Time: mustParseDate(t, "2024-02-01"), Factor: FactorPlus, Value: 5},
		}
		result, err := AccumulateSequenceByInterval(sequence, 0, 3, IntervalDays)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Fatalf("got %d buckets want 1 (both dates should share an epoch-aligned bucket)", len(result))
		}
		if result[0].Store != 15 {
			t.Fatalf("got store %v want 15", result[0].Store)
		}
	})

	t.Run("step=7 days: consistent bucketing across month boundary", func(t *testing.T) {
		sequence := []SeqFactor{
			{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-28"), Factor: FactorPlus, Value: 10},
			{ID: "2", Tag: "a", Time: mustParseDate(t, "2024-02-04"), Factor: FactorPlus, Value: 5},
		}
		// 2024-01-28 and 2024-02-04 are exactly 7 days apart — different buckets.
		result, err := AccumulateSequenceByInterval(sequence, 0, 7, IntervalDays)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("got %d buckets want 2", len(result))
		}
	})
}

func TestAccumulateSequenceByGrouping(t *testing.T) {
	sequence := []SeqFactor{
		{ID: "1", Tag: "a", Time: mustParseDate(t, "2024-01-01"), Factor: FactorPlus, Value: 10},
		{ID: "2", Tag: "a", Time: mustParseDate(t, "2024-04-01"), Factor: FactorMinus, Value: 3},
	}

	t.Run("produces same result as AccumulateSequenceByInterval", func(t *testing.T) {
		byInterval, err := AccumulateSequenceByInterval(sequence, 0, 3, IntervalMonths)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		byGrouping, err := AccumulateSequenceByGrouping(sequence, 0, TimeGrouping{Step: 3, IntervalType: IntervalMonths})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(byInterval) != len(byGrouping) {
			t.Fatalf("length mismatch: byInterval=%d byGrouping=%d", len(byInterval), len(byGrouping))
		}
		for i := range byInterval {
			if byInterval[i].Store != byGrouping[i].Store {
				t.Fatalf("store mismatch at [%d]", i)
			}
		}
	})

	t.Run("returns error for invalid grouping", func(t *testing.T) {
		_, err := AccumulateSequenceByGrouping(sequence, 0, TimeGrouping{Step: -1, IntervalType: IntervalDays})
		if err == nil {
			t.Fatalf("expected error for negative step")
		}
	})
}
