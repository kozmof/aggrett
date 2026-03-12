package aggrett

import "testing"

func TestAggregate(t *testing.T) {
	yesterday := mustParseDate(t, "2024-01-01")
	today := mustParseDate(t, "2024-01-02")
	tomorrow := mustParseDate(t, "2024-01-03")

	sequence := []SeqFactor{
		{ID: "1", Tag: "test", Time: today, Factor: FactorPlus, Value: 4},
		{ID: "2", Tag: "test", Time: today, Factor: FactorMinus, Value: 2},
		{ID: "3", Tag: "test", Time: yesterday, Factor: FactorPlus, Value: 10},
		{ID: "4", Tag: "test", Time: tomorrow, Factor: FactorMinus, Value: 3},
		{ID: "5", Tag: "other tag", Time: tomorrow, Factor: FactorMinus, Value: 5},
	}

	accum, err := Aggregate(sequence, 10, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(accum) != 3 {
		t.Fatalf("got %d items want 3", len(accum))
	}

	if accum[0].Store != 20 {
		t.Fatalf("accum[0].store got %v want 20", accum[0].Store)
	}
	if accum[0].Breakdown["test"].Delta != 10 {
		t.Fatalf("accum[0].breakdown[test].store got %v want 10", accum[0].Breakdown["test"].Delta)
	}
	if len(accum[0].Breakdown["test"].IDs) != 1 || accum[0].Breakdown["test"].IDs[0] != "3" {
		t.Fatalf("accum[0].breakdown[test].ids got %#v want [3]", accum[0].Breakdown["test"].IDs)
	}

	if accum[1].Store != 22 {
		t.Fatalf("accum[1].store got %v want 22", accum[1].Store)
	}
	if accum[1].Breakdown["test"].Delta != 2 {
		t.Fatalf("accum[1].breakdown[test].store got %v want 2", accum[1].Breakdown["test"].Delta)
	}
	if len(accum[1].Breakdown["test"].IDs) != 2 || accum[1].Breakdown["test"].IDs[0] != "1" || accum[1].Breakdown["test"].IDs[1] != "2" {
		t.Fatalf("accum[1].breakdown[test].ids got %#v want [1 2]", accum[1].Breakdown["test"].IDs)
	}

	if accum[2].Store != 14 {
		t.Fatalf("accum[2].store got %v want 14", accum[2].Store)
	}
	if accum[2].Breakdown["test"].Delta != -3 {
		t.Fatalf("accum[2].breakdown[test].store got %v want -3", accum[2].Breakdown["test"].Delta)
	}
	if accum[2].Breakdown["other tag"].Delta != -5 {
		t.Fatalf("accum[2].breakdown[other tag].store got %v want -5", accum[2].Breakdown["other tag"].Delta)
	}
	if len(accum[2].Breakdown["test"].IDs) != 1 || accum[2].Breakdown["test"].IDs[0] != "4" {
		t.Fatalf("accum[2].breakdown[test].ids got %#v want [4]", accum[2].Breakdown["test"].IDs)
	}
	if len(accum[2].Breakdown["other tag"].IDs) != 1 || accum[2].Breakdown["other tag"].IDs[0] != "5" {
		t.Fatalf("accum[2].breakdown[other tag].ids got %#v want [5]", accum[2].Breakdown["other tag"].IDs)
	}
}

func TestAggregateByInterval(t *testing.T) {
	sequence := []SeqFactor{
		{ID: "1", Tag: "rent", Time: mustParseDate(t, "2024-01-15"), Factor: FactorPlus, Value: 10},
		{ID: "2", Tag: "food", Time: mustParseDate(t, "2024-04-30"), Factor: FactorMinus, Value: 2},
		{ID: "3", Tag: "rent", Time: mustParseDate(t, "2024-05-01"), Factor: FactorPlus, Value: 7},
		{ID: "4", Tag: "food", Time: mustParseDate(t, "2024-08-31"), Factor: FactorPlus, Value: 1},
		{ID: "5", Tag: "rent", Time: mustParseDate(t, "2024-09-01"), Factor: FactorMinus, Value: 3},
	}

	accum, err := AggregateByInterval(sequence, 0, nil, 4, IntervalMonths)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(accum) != 3 {
		t.Fatalf("got %d items want 3", len(accum))
	}

	mustTimeEqual(t, accum[0].Time, mustParseDate(t, "2024-01-01"))
	mustTimeEqual(t, accum[1].Time, mustParseDate(t, "2024-05-01"))
	mustTimeEqual(t, accum[2].Time, mustParseDate(t, "2024-09-01"))

	if accum[0].Store != 8 || accum[1].Store != 16 || accum[2].Store != 13 {
		t.Fatalf("unexpected stores %#v", []float64{accum[0].Store, accum[1].Store, accum[2].Store})
	}
	if accum[0].Breakdown["rent"].Delta != 10 || accum[0].Breakdown["food"].Delta != -2 {
		t.Fatalf("unexpected first breakdown %#v", accum[0].Breakdown)
	}
	if accum[1].Breakdown["rent"].Delta != 7 || accum[1].Breakdown["food"].Delta != 1 {
		t.Fatalf("unexpected second breakdown %#v", accum[1].Breakdown)
	}
	if accum[2].Breakdown["rent"].Delta != -3 {
		t.Fatalf("unexpected third breakdown %#v", accum[2].Breakdown)
	}

	t.Run("returns error for non-positive step", func(t *testing.T) {
		_, err := AggregateByInterval(sequence, 0, nil, 0, IntervalMonths)
		if err == nil {
			t.Fatalf("expected error for step=0")
		}
	})

	t.Run("returns error for invalid interval type", func(t *testing.T) {
		_, err := AggregateByInterval(sequence, 0, nil, 1, IntervalType("bad"))
		if err == nil {
			t.Fatalf("expected error for invalid interval type")
		}
	})
}
