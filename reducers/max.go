package reducers

import (
	"cmp"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"github.com/djordje200179/meduce"
)

// NewMaxField creates a reducer that returns
// the value with maximal field.
//
// The reducer requires a getter function that
// returns fields that are compared.
// Comparator function is also required
// for comparing returned fields.
// Comparator function is also required
// for comparing returned fields.
func NewMaxField[KeyOut, ValueOut, ValueField any](
	getter func(value ValueOut) ValueField,
	comparator comparison.Comparator[ValueField],
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

// NewMaxOrderedField creates a reducer that returns
// the value with maximal field.
//
// The reducer requires a getter function that
// returns ordered fields that are natively compared.
func NewMaxOrderedField[KeyOut, ValueOut any, ValueField cmp.Ordered](getter func(value ValueOut) ValueField) meduce.Reducer[KeyOut, ValueOut] {
	return NewMaxField[KeyOut, ValueOut, ValueField](getter, cmp.Compare[ValueField])
}

// MaxPrimitive is a reducer that returns the maximum
// value of the values passed to it.
func MaxPrimitive[KeyOut any, ValueOut cmp.Ordered](_ KeyOut, values []ValueOut) ValueOut {
	maxValue := values[0]

	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue
}
