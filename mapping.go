package meduce

import (
	"runtime"
	"sync"
)

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) mapData() {
	threadsCount := runtime.NumCPU()

	var allMappersFinished sync.WaitGroup
	allMappersFinished.Add(threadsCount)

	keysArrays := make([][]KeyOut, threadsCount)
	valuesArrays := make([][]ValueOut, threadsCount)

	for i := 0; i < threadsCount; i++ {
		go mappingThread(
			process.KeyComparator,
			process.Mapper, process.Reducer,
			process.DataSource,
			&keysArrays[i], &valuesArrays[i],
			&allMappersFinished,
		)
	}

	allMappersFinished.Wait()

	process.mappedKeys, process.mappedValues = mergeMappedData(process.KeyComparator, keysArrays, valuesArrays)
}
