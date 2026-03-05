package aggrep

import "testing"

func TestAccumulate(t *testing.T) {
	t.Run("adds value with plus factor", func(t *testing.T) {
		if got := Accumulate(FactorPlus, 10, 5); got != 15 {
			t.Fatalf("got %v want 15", got)
		}
	})

	t.Run("subtracts value with minus factor", func(t *testing.T) {
		if got := Accumulate(FactorMinus, 10, 5); got != 5 {
			t.Fatalf("got %v want 5", got)
		}
	})

	t.Run("works with zero base value", func(t *testing.T) {
		if got := Accumulate(FactorPlus, 0, 100); got != 100 {
			t.Fatalf("got %v want 100", got)
		}
		if got := Accumulate(FactorMinus, 0, 100); got != -100 {
			t.Fatalf("got %v want -100", got)
		}
	})

	t.Run("works with zero value", func(t *testing.T) {
		if got := Accumulate(FactorPlus, 50, 0); got != 50 {
			t.Fatalf("got %v want 50", got)
		}
		if got := Accumulate(FactorMinus, 50, 0); got != 50 {
			t.Fatalf("got %v want 50", got)
		}
	})

	t.Run("handles negative results", func(t *testing.T) {
		if got := Accumulate(FactorMinus, 3, 10); got != -7 {
			t.Fatalf("got %v want -7", got)
		}
	})
}
