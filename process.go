package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc/functions"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"golang.org/x/exp/constraints"
	"sync"
)

type Process[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	KeyComparator functions.Comparator[KeyOut]

	Mapper    Mapper[KeyIn, ValueIn, KeyOut, ValueOut]
	Reducer   Reducer[KeyOut, ValueOut]
	Finalizer Finalizer[KeyOut, ValueOut]

	DataSource Source[KeyIn, ValueIn]

	mappedKeys   []KeyOut
	mappedValues []ValueOut

	collectingMutex sync.Mutex
	Collector       Collector[KeyOut, ValueOut]

	finishSignal sync.WaitGroup
}

func NewProcess[KeyIn, ValueIn, KeyOut, ValueOut any](
	keyComparator functions.Comparator[KeyOut],
	mapper Mapper[KeyIn, ValueIn, KeyOut, ValueOut], reducer Reducer[KeyOut, ValueOut], finalizer Finalizer[KeyOut, ValueOut],
	dataSource Source[KeyIn, ValueIn], collector Collector[KeyOut, ValueOut],
) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	//output = bufio.NewWriter(output)

	process := &Process[KeyIn, ValueIn, KeyOut, ValueOut]{
		KeyComparator: keyComparator,

		Mapper:    mapper,
		Reducer:   reducer,
		Finalizer: finalizer,

		DataSource: dataSource,

		Collector: collector,
	}

	return process
}

func NewProcessWithOrderedKeys[KeyIn, ValueIn any, KeyOut constraints.Ordered, ValueOut any](
	mapper Mapper[KeyIn, ValueIn, KeyOut, ValueOut], reducer Reducer[KeyOut, ValueOut], finalizer Finalizer[KeyOut, ValueOut],
	dataSource Source[KeyIn, ValueIn], collector Collector[KeyOut, ValueOut],
) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	return NewProcess(comparison.Ascending[KeyOut], mapper, reducer, finalizer, dataSource, collector)
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) Run() {
	process.mapData()
	process.reduceData()

	process.Collector.Finalize()
	process.finishSignal.Done()
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) WaitToFinish() Collector[KeyOut, ValueOut] {
	process.finishSignal.Add(1)
	process.finishSignal.Wait()

	return process.Collector
}
