package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc/functions"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"golang.org/x/exp/constraints"
	"log"
	"sync"
)

type Config[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	KeyComparator   functions.Comparator[KeyOut]
	ValueComparator functions.Comparator[ValueOut]

	Mapper    Mapper[KeyIn, ValueIn, KeyOut, ValueOut]
	Reducer   Reducer[KeyOut, ValueOut]
	Finalizer Finalizer[KeyOut, ValueOut]
	Filter    Filter[KeyOut, ValueOut]

	Source    Source[KeyIn, ValueIn]
	Collector Collector[KeyOut, ValueOut]
}

var nextUid = 0

type Process[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	uid int

	Config[KeyIn, ValueIn, KeyOut, ValueOut]

	mappedKeys   []KeyOut
	mappedValues []ValueOut

	collectingMutex sync.Mutex

	finishSignal sync.WaitGroup
}

func NewProcess[KeyIn, ValueIn, KeyOut, ValueOut any](config Config[KeyIn, ValueIn, KeyOut, ValueOut]) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	//output = bufio.NewWriter(output)

	if config.KeyComparator == nil {
		panic("KeyComparator must be set")
	}

	if config.Mapper == nil {
		panic("Mapper must be set")
	}

	if config.Reducer == nil {
		panic("Reducer must be set")
	}

	process := &Process[KeyIn, ValueIn, KeyOut, ValueOut]{
		uid: nextUid,

		Config: config,
	}

	nextUid++

	process.finishSignal.Add(1)

	return process
}

func NewDefaultProcess[KeyIn, ValueIn any, KeyOut constraints.Ordered, ValueOut any](
	config Config[KeyIn, ValueIn, KeyOut, ValueOut],
) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	config.KeyComparator = comparison.Ascending[KeyOut]

	return NewProcess(config)
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
