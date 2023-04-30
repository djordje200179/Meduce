package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc"
	"reflect"
	"runtime"
	"sync"
)

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) reduceData() {
	threadsCount := runtime.NumCPU()

	var barrier sync.WaitGroup
	barrier.Add(threadsCount)

	readyDataPool := make(chan reducingDataGroup[KeyOut, ValueOut], 100*threadsCount)

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

	for i := 0; i < threadsCount; i++ {
		go reducingThread(process, readyDataPool, &barrier)
	}

	if process.Logger != nil {
		process.Logger.Printf("Process %d: %d reducing threads were started\n", process.uid, threadsCount)
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
