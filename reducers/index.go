package reducers

// First is a reducer that returns the first
// value of the values passed to it.
func First[KeyOut, ValueOut any](_ KeyOut, values []ValueOut) ValueOut {
	return values[0]
}

// Last is a reducer that returns the last
// value of the values passed to it.
func Last[KeyOut, ValueOut any](_ KeyOut, values []ValueOut) ValueOut {
	return values[len(values)-1]
}
