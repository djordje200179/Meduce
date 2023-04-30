// Package meduce implements an interface to
// run MapReduce tasks on a single machine.
package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc"
	"github.com/djordje200179/extendedlibrary/misc/functions"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"golang.org/x/exp/constraints"
	"log"
	"sync"
	"sync/atomic"
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
//
// Zero value of Process has no configuration set and has invalid uid.
type Process[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	uid int

	Config[KeyIn, ValueIn, KeyOut, ValueOut]

	mappedKeys   []KeyOut
	mappedValues []ValueOut

	collectingMutex sync.Mutex
	linkBuffer      chan misc.Pair[KeyOut, ValueOut]

	mappingsFinished   atomic.Uint64
	reductionsFinished atomic.Uint64

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
func NewDefaultProcess[KeyIn, ValueIn any, KeyOut constraints.Ordered, ValueOut any](
	config Config[KeyIn, ValueIn, KeyOut, ValueOut],
) *Process[KeyIn, ValueIn, KeyOut, ValueOut] {
	config.KeyComparator = comparison.Ascending[KeyOut]

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
	if process.Logger != nil {
		process.Logger.Printf(
			"Process %d: mapping phase finished, %d mappings and %d reductions were made\n",
			process.uid, process.MappingsFinished(), process.reductionsFinished.Load())
	}

	process.reduceData()
	if process.Logger != nil {
		process.Logger.Printf("Process %d: reductions finished\n", process.uid)
	}

	process.processFinished.Done()
}

// WaitToFinish blocks until the MapReduce task is finished.
func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) WaitToFinish() Collector[KeyOut, ValueOut] {
	process.processFinished.Wait()

	return process.Collector
}

// MappingsFinished returns the number of mappings that
// have been finished so far.
func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) MappingsFinished() uint {
	return uint(process.mappingsFinished.Load())
}
