package meduce

import (
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

	for i := 0; i < threadsCount; i++ {
		go reducingThread(
			&process.Config,
			readyDataPool,
			process.collectorWrapper, &barrier,
		)
	}

	barrier.Wait()
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) collectorWrapper(key KeyOut, value ValueOut) {
	process.collectingMutex.Lock()
	process.Collector.Collect(key, value)
	process.collectingMutex.Unlock()
}
