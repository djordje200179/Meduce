package collectors

// SingleValueCollector is a collector that collects a single value.
//
// Zero value of SingleValueCollector is a valid collector
// that collects a value for the zero value of the key type.
type SingleValueCollector[KeyOut comparable, ValueOut any] struct {
	set bool

	wantedKey KeyOut
	value     ValueOut
}

// NewSingleValueCollector creates a new SingleValueCollector
// that collects a value for the given key.
func NewSingleValueCollector[KeyOut comparable, ValueOut any](wantedKey KeyOut) *SingleValueCollector[KeyOut, ValueOut] {
	collector := &SingleValueCollector[KeyOut, ValueOut]{
		wantedKey: wantedKey,
	}

	return collector
}

func (collector *SingleValueCollector[KeyOut, ValueOut]) Init() {

}

func (collector *SingleValueCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	if key == collector.wantedKey {
		collector.value = value
		collector.set = true
	}
}

func (collector *SingleValueCollector[KeyOut, ValueOut]) Finalize() {

}

// TryGet returns the collected value and true if the value was collected.
// If the value was not collected, it returns the default value and false.
func (collector *SingleValueCollector[KeyOut, ValueOut]) TryGet() (ValueOut, bool) {
	if collector.set {
		return collector.value, true
	} else {
		return collector.value, false
	}
}

// Get returns the collected value.
// If the value was not collected, it panics.
func (collector *SingleValueCollector[KeyOut, ValueOut]) Get() ValueOut {
	if collector.set {
		return collector.value
	} else {
		panic("Value was not set")
	}
}
