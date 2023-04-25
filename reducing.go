package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"runtime"
	"sync"
)

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) reduceData() {
	threadsCount := runtime.NumCPU()

	var barrier sync.WaitGroup
	barrier.Add(threadsCount)

	readyDataPool := make(chan reducingGroupData[KeyOut, ValueOut], 100*threadsCount)

	go func() {
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

			if firstIndex == lastIndex {
				value := process.mappedValues[firstIndex]

				barrier.Add(1)
				go func() {
					if process.finalizer != nil {
						process.finalizer(lastKey, &value)
					}
					process.collectData(lastKey, value)

					barrier.Done()
				}()

				continue
			}

			validValues := process.mappedValues[firstIndex : lastIndex+1]

			reducerData := reducingGroupData[KeyOut, ValueOut]{
				key:    lastKey,
				values: validValues,
			}

			readyDataPool <- reducerData
		}

		close(readyDataPool)
	}()

	for i := 0; i < threadsCount; i++ {
		go reduceData(
			process.reducer, process.finalizer,
			readyDataPool,
			process.collectData, &barrier,
		)
	}

	barrier.Wait()
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) collectData(key KeyOut, value ValueOut) {
	process.collectingMutex.Lock()
	process.collector.Collect(key, value)
	process.collectingMutex.Unlock()
}
