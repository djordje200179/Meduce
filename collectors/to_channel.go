package collectors

import "github.com/djordje200179/extendedlibrary/misc"

// ChannelCollector is a collector that collects key-value pairs into a channel.
type ChannelCollector[KeyOut, ValueOut any] chan<- misc.Pair[KeyOut, ValueOut]

// NewChannelCollector creates a new ChannelCollector
// with the specified buffer size.
func NewChannelCollector[KeyOut, ValueOut any](bufferSize int) ChannelCollector[KeyOut, ValueOut] {
	return make(chan misc.Pair[KeyOut, ValueOut], bufferSize)
}

func (collector ChannelCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	collector <- misc.Pair[KeyOut, ValueOut]{key, value}
}

func (collector ChannelCollector[KeyOut, ValueOut]) Finalize() {
	close(collector)
}
