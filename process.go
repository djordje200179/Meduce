package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc/functions"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"golang.org/x/exp/constraints"
	"log"
	"sync"
)

var nextUid int = 0

type Process[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	uid int

	KeyComparator   functions.Comparator[KeyOut]
	ValueComparator functions.Comparator[ValueOut]

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
	keyComparator functions.Comparator[KeyOut], valueComparator functions.Comparator[ValueOut],
	mapper Mapper[KeyIn, ValueIn, KeyOut, ValueOut], reducer Reducer[KeyOut, ValueOut], finalizer Finalizer[KeyOut, ValueOut],
	dataSource Source[KeyIn, ValueIn], collector Collector[KeyOut, ValueOut],
) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	//output = bufio.NewWriter(output)

	process := &Process[KeyIn, ValueIn, KeyOut, ValueOut]{
		uid: nextUid,

		KeyComparator:   keyComparator,
		ValueComparator: valueComparator,

		Mapper:    mapper,
		Reducer:   reducer,
		Finalizer: finalizer,

		DataSource: dataSource,

		Collector: collector,
	}

	nextUid++

	process.finishSignal.Add(1)

	return process
}

func NewProcessWithOrderedKeys[KeyIn, ValueIn any, KeyOut constraints.Ordered, ValueOut any](
	mapper Mapper[KeyIn, ValueIn, KeyOut, ValueOut], reducer Reducer[KeyOut, ValueOut], finalizer Finalizer[KeyOut, ValueOut],
	dataSource Source[KeyIn, ValueIn], collector Collector[KeyOut, ValueOut],
) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	return NewProcess(comparison.Ascending[KeyOut], nil, mapper, reducer, finalizer, dataSource, collector)
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) Run(verbose bool) {
	if verbose {
		log.Printf("Process %d: started", process.uid)
	}
	process.mapData()
	if verbose {
		log.Printf("Process %d: mappings finished", process.uid)
	}
	process.reduceData()
	if verbose {
		log.Printf("Process %d: reductions finished", process.uid)
	}

	process.Collector.Finalize()
	process.finishSignal.Done()
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) WaitToFinish() Collector[KeyOut, ValueOut] {
	process.finishSignal.Wait()

	return process.Collector
}
