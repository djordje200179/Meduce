package reducers

import "cmp"

// SumPrimitive is a reducer that returns the sum
// of the values passed to it.
func SumPrimitive[KeyOut any, ValueOut cmp.Ordered](_ KeyOut, values []ValueOut) ValueOut {
	var sum ValueOut

	for _, value := range values {
		sum += value
	}

	return sum
}
