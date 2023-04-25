package meduce

import "github.com/djordje200179/extendedlibrary/misc"

type Source[KeyIn, ValueIn any] <-chan misc.Pair[KeyIn, ValueIn]
type Emitter[KeyOut, ValueOut any] func(key KeyOut, value ValueOut)
type Mapper[KeyIn, ValueIn, KeyOut, ValueOut any] func(key KeyIn, value ValueIn, emit Emitter[KeyOut, ValueOut])
type Reducer[KeyOut, ValueOut any] func(key KeyOut, values []ValueOut) ValueOut
type Finalizer[KeyOut, ValueOut any] func(key KeyOut, valueRef *ValueOut)

type Collector[KeyOut, ValueOut any] interface {
	Collect(key KeyOut, value ValueOut)
	Finalize()
}
