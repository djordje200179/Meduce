package meduce

import "github.com/djordje200179/extendedlibrary/misc"

// A Source is a channel that is created by user
// and from which key-value pairs are read.
//
// Preferably, it should be buffered to avoid
// blocking and context switches.
type Source[KeyIn, ValueIn any] <-chan misc.Pair[KeyIn, ValueIn]

// An Emitter is a function that is supplied by library.
//
// It is passed to user's Mapper function,
// and is called to emit key-value pairs.
type Emitter[KeyOut, ValueOut any] func(key KeyOut, value ValueOut)

// A Mapper is a function that is created by user
// and is used to map input data to key-value pairs.
type Mapper[KeyIn, ValueIn, KeyOut, ValueOut any] func(key KeyIn, value ValueIn, emit Emitter[KeyOut, ValueOut])

// A Reducer is a function that is created by user
// and is used to reduce values to single value.
//
// It is called for all keys multiple times,
// until all values for that key are reduced.
// It should be idempotent and have no side effects.
type Reducer[KeyOut, ValueOut any] func(key KeyOut, values []ValueOut) ValueOut

// A Finalizer is a function that is created by user
// and is used to finalize key-value pairs.
//
// It is called after all values for a key were reduced
// to a single value.
type Finalizer[KeyOut, ValueOut any] func(key KeyOut, valueRef *ValueOut)

// A Filter is a function that is created by user
// and is used to filter processed key-value pairs.
type Filter[KeyOut, ValueOut any] func(key KeyOut, valueRef *ValueOut) bool

// A Collector is an entity that is supplied by user
// and is used to collect processed key-value pairs.
type Collector[KeyOut, ValueOut any] interface {
	Init()                              // Init is called just before collecting starts
	Collect(key KeyOut, value ValueOut) // Collect is called for each processed key-value pair
	Finalize()                          // Finalize is called after all key-value pairs were processed
}
