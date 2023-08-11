package reducers

import (
	"cmp"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"github.com/djordje200179/meduce"
)

// NewMinField creates a reducer that returns
// the value with minimal field.
//
// The reducer requires a getter function that
// returns fields that are compared.
// Comparator function is also required
// for comparing returned fields.
func NewMinField[KeyOut, ValueOut, ValueField any](
	getter func(value ValueOut) ValueField,
	comparator comparison.Comparator[ValueField],
) meduce.Reducer[KeyOut, ValueOut] {
	return func(_ KeyOut, values []ValueOut) ValueOut {
		minValue := values[0]
		minField := getter(values[0])

		for _, value := range values {
			field := getter(value)
			if comparator(field, minField) == comparison.FirstSmaller {
				minField = field
				minValue = value
			}
		}

		return minValue
	}
}

// NewMinOrderedField creates a reducer that returns
// the value with minimal field.
//
// The reducer requires a getter function that
// returns ordered fields that are natively compared.
func NewMinOrderedField[KeyOut, ValueOut any, ValueField cmp.Ordered](getter func(value ValueOut) ValueField) meduce.Reducer[KeyOut, ValueOut] {
	return NewMinField[KeyOut, ValueOut, ValueField](getter, cmp.Compare[ValueField])
}

// MinPrimitive is a reducer that returns the minimum
// value of the values passed to it.
func MinPrimitive[KeyOut any, ValueOut cmp.Ordered](_ KeyOut, values []ValueOut) ValueOut {
	minValue := values[0]

	for _, value := range values {
		if value < minValue {
			minValue = value
		}
	}

	return minValue
}
