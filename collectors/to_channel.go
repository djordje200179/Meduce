package collectors

import "github.com/djordje200179/extendedlibrary/misc"

type ChannelCollector[KeyOut, ValueOut any] chan<- misc.Pair[KeyOut, ValueOut]

func NewChannelCollector[KeyOut, ValueOut any](bufferSize int) ChannelCollector[KeyOut, ValueOut] {
	return make(chan misc.Pair[KeyOut, ValueOut], bufferSize)
}

func (collector ChannelCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	collector <- misc.Pair[KeyOut, ValueOut]{key, value}
}

func (collector ChannelCollector[KeyOut, ValueOut]) Finalize() {
	close(collector)
}
