package collectors

type SingleValueCollector[KeyOut comparable, ValueOut any] struct {
	set bool

	wantedKey KeyOut
	value     ValueOut
}

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

func (collector *SingleValueCollector[KeyOut, ValueOut]) TryGet() (ValueOut, bool) {
	if collector.set {
		return collector.value, true
	} else {
		return collector.value, false
	}
}

func (collector *SingleValueCollector[KeyOut, ValueOut]) Get() ValueOut {
	if collector.set {
		return collector.value
	} else {
		panic("Value was not set")
	}
}
