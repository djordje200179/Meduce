package collectors

import "github.com/djordje200179/meduce"

type MapCollector[KeyOut comparable, ValueOut any] struct {
	m map[KeyOut]ValueOut
}

func NewMapCollector[KeyOut comparable, ValueOut any]() meduce.Collector[KeyOut, ValueOut] {
	collector := MapCollector[KeyOut, ValueOut]{
		m: make(map[KeyOut]ValueOut),
	}

	return collector
}

func (collector MapCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	collector.m[key] = value
}

func (collector MapCollector[KeyOut, ValueOut]) Finalize() {

}
