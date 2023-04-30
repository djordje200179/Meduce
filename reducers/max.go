package reducers

import (
	"github.com/djordje200179/extendedlibrary/misc/functions"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"github.com/djordje200179/meduce"
	"golang.org/x/exp/constraints"
)

// Max returns a reducer that returns the maximum
// value of the values passed to it.
//
// The reducer requires a getter function that
// returns the field that should be compared
// by the comparator function.
func Max[KeyOut, ValueOut, ValueField any](
	getter func(value ValueOut) ValueField,
	comparator functions.Comparator[ValueField],
) meduce.Reducer[KeyOut, ValueOut] {
	return func(_ KeyOut, values []ValueOut) ValueOut {
		maxValue := values[0]
		maxField := getter(values[0])

		for _, value := range values {
			field := getter(value)
			if comparator(field, maxField) == comparison.FirstBigger {
				maxField = field
				maxValue = value
			}
		}

		return maxValue
	}
}

// MaxOrdered returns a reducer that returns the maximum
// value of the values passed to it.
//
// The reducer requires a getter function that
// returns the field that should be compared.
func MaxOrdered[KeyOut, ValueOut any, ValueField constraints.Ordered](getter func(value ValueOut) ValueField) meduce.Reducer[KeyOut, ValueOut] {
	return Max[KeyOut, ValueOut, ValueField](getter, comparison.Compare[ValueField])
}

// MaxPrimitive is a reducer that returns the maximum
// value of the values passed to it.
func MaxPrimitive[KeyOut any, ValueOut constraints.Ordered](_ KeyOut, values []ValueOut) ValueOut {
	maxValue := values[0]

	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue
}
