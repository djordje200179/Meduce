package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"sync"
)

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) reduceData() {
	var barrier sync.WaitGroup

	lastIndex := -1
	for i := 1; i <= len(process.mappedKeys); i++ {
		lastKey := process.mappedKeys[i-1]

		if i != len(process.mappedKeys) {
			currentKey := process.mappedKeys[i]

			if process.keyComparator(lastKey, currentKey) == comparison.Equal {
				continue
			}
		}

		firstIndex := lastIndex + 1
		lastIndex = i - 1

		barrier.Add(1)

		if firstIndex == lastIndex {
			value := process.mappedValues[firstIndex]
			go writeOnlyData(process.collectData, &barrier, lastKey, value)

			continue
		}

		validValues := process.mappedValues[firstIndex : lastIndex+1]
		go reduceData(
			process.reducer, process.finalizer,
			process.collectData, &barrier,
			lastKey, validValues,
		)
	}

	barrier.Wait()
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) collectData(key KeyOut, value ValueOut) {
	process.collectingMutex.Lock()
	process.collector.Collect(key, value)
	process.collectingMutex.Unlock()
}
