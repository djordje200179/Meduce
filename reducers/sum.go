package reducers

import "golang.org/x/exp/constraints"

// Sum is a reducer that returns the sum
// of the values passed to it.
func SumPrimitive[KeyOut any, ValueOut constraints.Ordered](_ KeyOut, values []ValueOut) ValueOut {
	var sum ValueOut

	for _, value := range values {
		sum += value
	}

	return sum
}
