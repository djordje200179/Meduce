// Package meduce implements an interface to
// run MapReduce tasks on a single machine.
package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc/functions"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"golang.org/x/exp/constraints"
	"log"
	"sync"
)

// A Config is a configuration for a single MapReduce task.
type Config[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	// KeyComparator and ValueComparator are used to sort key-value pairs
	// before they are passed to the Reducer.
	// KeyComparator is used as primary comparator,
	// and ValueComparator is used as secondary.
	KeyComparator   functions.Comparator[KeyOut]
	ValueComparator functions.Comparator[ValueOut]

	Mapper    Mapper[KeyIn, ValueIn, KeyOut, ValueOut]
	Reducer   Reducer[KeyOut, ValueOut]
	Finalizer Finalizer[KeyOut, ValueOut]
	Filter    Filter[KeyOut, ValueOut]

	Source    Source[KeyIn, ValueIn]
	Collector Collector[KeyOut, ValueOut]

	Logger *log.Logger
}

var nextUid = 0

// A Process is an instance of a single MapReduce task.
type Process[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	uid int

	Config[KeyIn, ValueIn, KeyOut, ValueOut]

	mappedKeys   []KeyOut
	mappedValues []ValueOut

	collectingMutex sync.Mutex

	finishSignal sync.WaitGroup
}

// NewProcess creates a new Process with given configuration.
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

// NewDefaultProcess creates a new Process with default key comparator for ordered keys.
func NewDefaultProcess[KeyIn, ValueIn any, KeyOut constraints.Ordered, ValueOut any](
	config Config[KeyIn, ValueIn, KeyOut, ValueOut],
) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	config.KeyComparator = comparison.Ascending[KeyOut]

	return NewProcess(config)
}

// Run starts the MapReduce task and blocks until it is finished.
//
// If logger is set, it will be used to log the progress.
func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) Run() {
	if process.Logger != nil {
		process.Logger.Printf("Process %d: started", process.uid)
	}
	process.mapData()
	if process.Logger != nil {
		process.Logger.Printf("Process %d: mappings finished", process.uid)
	}
	process.reduceData()
	if process.Logger != nil {
		process.Logger.Printf("Process %d: reductions finished", process.uid)
	}

	process.Collector.Finalize()
	process.finishSignal.Done()
}

// WaitToFinish blocks until the MapReduce task is finished.
func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) WaitToFinish() Collector[KeyOut, ValueOut] {
	process.finishSignal.Wait()

	return process.Collector
}
