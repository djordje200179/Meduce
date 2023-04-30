package collectors

// SingleValueCollector is a collector that collects a single value.
//
// Zero value of SingleValueCollector is a valid collector.
type SingleValueCollector[KeyOut comparable, ValueOut any] struct {
	set bool

	key   KeyOut
	value ValueOut
}

// NewSingleValueCollector creates a new SingleValueCollector
// that collects a value for the given key.
func NewSingleValueCollector[KeyOut comparable, ValueOut any]() *SingleValueCollector[KeyOut, ValueOut] {
	return &SingleValueCollector[KeyOut, ValueOut]{}
}

func (collector *SingleValueCollector[KeyOut, ValueOut]) Init() {

}

func (collector *SingleValueCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	if collector.set {
		panic("Already collected a value")
	}

	collector.set = true
	collector.key = key
	collector.value = value
}

func (collector *SingleValueCollector[KeyOut, ValueOut]) Finalize() {

}

// Get returns the collected key-value pair.
// If no key-value pair was collected, panics.
func (collector *SingleValueCollector[KeyOut, ValueOut]) Get() (KeyOut, ValueOut) {
	if !collector.set {
		panic("No value collected")
	}

	return collector.key, collector.value
}

// Key returns the collected key.
// If no key-value pair was collected, panics.
func (collector *SingleValueCollector[KeyOut, ValueOut]) Key() KeyOut {
	if !collector.set {
		panic("No value collected")
	}

	return collector.key
}

// Value returns the collected value.
// If no key-value pair was collected, panics.
func (collector *SingleValueCollector[KeyOut, ValueOut]) Value() ValueOut {
	if !collector.set {
		panic("No value collected")
	}

	return collector.value
}

// IsSet returns true if a key-value pair was collected.
func (collector *SingleValueCollector[KeyOut, ValueOut]) IsSet() bool {
	return collector.set
}
