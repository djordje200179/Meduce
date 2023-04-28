package meduce

import "github.com/djordje200179/extendedlibrary/misc"

// A Source is a channel that is supplied by user
// and is used to read input data.
// Preferably, it should be buffered to avoid too
// many context switches and blocking.
type Source[KeyIn, ValueIn any] <-chan misc.Pair[KeyIn, ValueIn]

// An Emitter is a function that is supplied by library
// which user calls to emit key-value pairs.
type Emitter[KeyOut, ValueOut any] func(key KeyOut, value ValueOut)

// A Mapper is a function that is supplied by user
// and is used to map input data to key-value pairs.
type Mapper[KeyIn, ValueIn, KeyOut, ValueOut any] func(key KeyIn, value ValueIn, emit Emitter[KeyOut, ValueOut])

// A Reducer is a function that is supplied by user
// and is used to reduce key-value pairs to single value.
type Reducer[KeyOut, ValueOut any] func(key KeyOut, values []ValueOut) ValueOut

// A Finalizer is a function that is supplied by user
// and is used to finalize key-value pairs.
type Finalizer[KeyOut, ValueOut any] func(key KeyOut, valueRef *ValueOut)

// A Filter is a function that is supplied by user
// and is used to filter processed key-value pairs.
type Filter[KeyOut, ValueOut any] func(key KeyOut, valueRef *ValueOut) bool

// A Collector is a function that is supplied by user
// and is used to collect processed key-value pairs.
type Collector[KeyOut, ValueOut any] interface {
	Collect(key KeyOut, value ValueOut) // Collect is called for each processed key-value pair
	Finalize()                          // Finalize is called after all key-value pairs are processed
}
