package collectors

import "github.com/djordje200179/meduce"

type MapCollector[KeyOut comparable, ValueOut any] map[KeyOut]ValueOut

func NewMapCollector[KeyOut comparable, ValueOut any]() meduce.Collector[KeyOut, ValueOut] {
	return make(MapCollector[KeyOut, ValueOut])
}

func (collector MapCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	collector[key] = value
}

func (collector MapCollector[KeyOut, ValueOut]) Finalize() {

}
