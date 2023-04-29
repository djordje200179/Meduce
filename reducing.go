package meduce

import (
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

	process.Collector.Init()

	for i := 0; i < threadsCount; i++ {
		go reducingThread(
			&process.Config,
			readyDataPool,
			process.collectorWrapper, &barrier,
		)
	}

	barrier.Wait()

	process.Collector.Finalize()
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) collectorWrapper(key KeyOut, value ValueOut) {
	if reflect.TypeOf(process.Collector).Kind() != reflect.Chan {
		process.collectingMutex.Lock()
		defer process.collectingMutex.Unlock()
	}

	process.Collector.Collect(key, value)
}
