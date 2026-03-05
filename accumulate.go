package aggrep

import "fmt"

// Accumulate applies a factor to a previous value.
func Accumulate(factor Factor, prevValue, value float64) float64 {
	switch factor {
	case FactorPlus:
		return prevValue + value
	case FactorMinus:
		return prevValue - value
	default:
		panic(fmt.Sprintf("unknown factor type: %q", factor))
	}
}
