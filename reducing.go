package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc"
	"reflect"
	"runtime"
	"sync"
)

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) reduceData() {
	groupsCount := process.estimateGroupsCount()

	var threadsCount int
	if groupsCount > runtime.NumCPU() {
		threadsCount = runtime.NumCPU()
	} else {
		threadsCount = groupsCount
	}

	readyDataPool := make(chan reducingDataGroup[KeyOut, ValueOut], groupsCount)
	go reducingDataGenerationThread(
		process.KeyComparator,
		process.mappedKeys, process.mappedValues,
		readyDataPool,
	)

	if process.Collector != nil {
		process.Collector.Init()
		defer process.Collector.Finalize()
	} else {
		go process.runNext()
		defer close(process.linkBuffer)
	}

	var barrier sync.WaitGroup
	barrier.Add(threadsCount)

	process.reducingThreads = make([]reducingThread[KeyIn, ValueIn, KeyOut, ValueOut], threadsCount)
	for i := range process.reducingThreads {
		process.reducingThreads[i].Process = process

		go process.reducingThreads[i].run(readyDataPool, &barrier)
	}

	if process.Logger != nil {
		var message string
		if threadsCount == 1 {
			message = "Process %d: %d reducing thread was started\n"
		} else {
			message = "Process %d: %d reducing threads were started\n"
		}

		process.Logger.Printf(message, process.uid, threadsCount)
	}

	barrier.Wait()

	if process.Logger != nil {
		process.Logger.Printf("Process %d: all reducing threads finished\n", process.uid)
	}
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) collect(key KeyOut, value ValueOut) {
	if process.Collector == nil {
		process.linkBuffer <- misc.Pair[KeyOut, ValueOut]{key, value}
		return
	}

	if reflect.TypeOf(process.Collector).Kind() != reflect.Chan {
		process.collectingMutex.Lock()
		defer process.collectingMutex.Unlock()
	}

	process.Collector.Collect(key, value)
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) estimateGroupsCount() int {
	var combinationsCount int

	for _, thread := range process.mappingThreads {
		if thread.combinationsCount > combinationsCount {
			combinationsCount = thread.combinationsCount
		}
	}

	return combinationsCount
}
