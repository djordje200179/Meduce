package meduce

import (
	"runtime"
	"sync"
)

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) mapData() {
	threadsCount := runtime.NumCPU()

	var barrier sync.WaitGroup
	barrier.Add(threadsCount)

	keysArrays := make([][]KeyOut, threadsCount)
	valuesArrays := make([][]ValueOut, threadsCount)

	for i := 0; i < threadsCount; i++ {
		go mapData(
			process.keyComparator,
			process.mapper, process.reducer,
			process.dataSource,
			&keysArrays[i], &valuesArrays[i],
			&barrier,
		)
	}

	barrier.Wait()

	process.mappedKeys, process.mappedValues = mergeMappedData(process.keyComparator, keysArrays, valuesArrays)
}
