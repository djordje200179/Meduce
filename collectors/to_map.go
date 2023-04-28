package collectors

import "github.com/djordje200179/meduce"

// MapCollector is a collector that collects key-value pairs into a map.
type MapCollector[KeyOut comparable, ValueOut any] map[KeyOut]ValueOut

// NewMapCollector creates a new MapCollector.
func NewMapCollector[KeyOut comparable, ValueOut any]() meduce.Collector[KeyOut, ValueOut] {
	return make(MapCollector[KeyOut, ValueOut])
}

func (collector MapCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	collector[key] = value
}

func (collector MapCollector[KeyOut, ValueOut]) Finalize() {

}
