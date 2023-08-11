// Package meduce implements an interface to
// run MapReduce tasks on a single machine.
package meduce

import (
	"cmp"
	"github.com/djordje200179/extendedlibrary/misc"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"log"
	"sync"
)

// A Config is a configuration for a single MapReduce task.
type Config[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	// KeyComparator and ValueComparator are used to sort key-value pairs
	// before they are passed to the Reducer.
	// KeyComparator is used as primary comparator,
	// and ValueComparator is used as secondary.
	KeyComparator   comparison.Comparator[KeyOut]
	ValueComparator comparison.Comparator[ValueOut]

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
//
// Zero value of Process has no configuration set and has invalid uid.
type Process[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	uid int

	Config[KeyIn, ValueIn, KeyOut, ValueOut]

	mappingThreads  []mappingThread[KeyIn, ValueIn, KeyOut, ValueOut]
	reducingThreads []reducingThread[KeyIn, ValueIn, KeyOut, ValueOut]

	mappedKeys   []KeyOut
	mappedValues []ValueOut

	collectingMutex sync.Mutex
	linkBuffer      chan misc.Pair[KeyOut, ValueOut]

	processFinished sync.WaitGroup

	runNext func()
}

// NewProcess creates a new Process with given configuration.
func NewProcess[KeyIn, ValueIn, KeyOut, ValueOut any](config Config[KeyIn, ValueIn, KeyOut, ValueOut]) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	nextUid++

	process := &Process[KeyIn, ValueIn, KeyOut, ValueOut]{
		uid: nextUid,

		Config: config,
	}

	process.processFinished.Add(1)

	return process
}

// NewDefaultProcess creates a new Process with default key comparator for ordered keys.
func NewDefaultProcess[KeyIn, ValueIn any, KeyOut cmp.Ordered, ValueOut any](
	config Config[KeyIn, ValueIn, KeyOut, ValueOut],
) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	config.KeyComparator = cmp.Compare[KeyOut]

	return NewProcess(config)
}

// Link links two processes together.
func Link[KeyOld, ValueOld, KeyIn, ValueIn, KeyOut, ValueOut any](
	prevProcess *Process[KeyOld, ValueOld, KeyIn, ValueIn],
	nextProcess *Process[KeyIn, ValueIn, KeyOut, ValueOut],
) {
	LinkWithBufferSize(prevProcess, nextProcess, 100)
}

// LinkWithBufferSize links two processes together with a buffer of given size.
//
// bufferSize is the size of the buffer that will be created to link the processes.
func LinkWithBufferSize[KeyOld, ValueOld, KeyIn, ValueIn, KeyOut, ValueOut any](
	prevProcess *Process[KeyOld, ValueOld, KeyIn, ValueIn],
	nextProcess *Process[KeyIn, ValueIn, KeyOut, ValueOut],
	bufferSize int,
) {
	buffer := make(chan misc.Pair[KeyIn, ValueIn], bufferSize)

	prevProcess.linkBuffer = buffer
	nextProcess.Source = buffer

	prevProcess.runNext = nextProcess.Run
}

// Run starts the MapReduce task and blocks until it is finished.
//
// If logger is set, it will be used to log the progress.
func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) Run() {
	if process.KeyComparator == nil {
		panic("KeyComparator must be set")
	}

	if process.Mapper == nil {
		panic("Mapper must be set")
	}

	if process.Reducer == nil {
		panic("Reducer must be set")
	}

	if process.Source == nil {
		panic("Source must be set")
	}

	if process.Collector == nil && process.linkBuffer == nil {
		panic("Collector must be set")
	}

	if process.Logger != nil {
		process.Logger.Printf("Process %d: started\n", process.uid)
	}

	process.mapData()
	process.reduceData()

	process.processFinished.Done()
}

// WaitToFinish blocks until the MapReduce task is finished.
func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) WaitToFinish() Collector[KeyOut, ValueOut] {
	process.processFinished.Wait()

	return process.Collector
}
