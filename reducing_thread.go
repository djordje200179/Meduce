package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc/functions"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"sync"
)

type reducingDataGroup[KeyOut, ValueOut any] struct {
	key    KeyOut
	values []ValueOut
}

func reducingThread[KeyIn, ValueIn, KeyOut, ValueOut any](
	process *Process[KeyIn, ValueIn, KeyOut, ValueOut],
	dataPool <-chan reducingDataGroup[KeyOut, ValueOut],
	finishSignal *sync.WaitGroup,
) {
	for groupData := range dataPool {
		var reducedValue ValueOut
		if len(groupData.values) == 1 {
			reducedValue = groupData.values[0]
		} else {
			reducedValue = process.Reducer(groupData.key, groupData.values)
		}

		if process.Finalizer != nil {
			process.Finalizer(groupData.key, &reducedValue)
		}

		if process.Filter == nil || process.Filter(groupData.key, &reducedValue) {
			process.collect(groupData.key, reducedValue)
		}
	}

	finishSignal.Done()
}

func reducingDataGenerationThread[KeyOut, ValueOut any](
	keyComparator functions.Comparator[KeyOut],
	mappedKeys []KeyOut, mappedValues []ValueOut,
	readyDataPool chan<- reducingDataGroup[KeyOut, ValueOut],
) {
	lastIndex := -1
	for i := 1; i <= len(mappedKeys); i++ {
		lastKey := mappedKeys[i-1]

		if i != len(mappedKeys) {
			currentKey := mappedKeys[i]

			if keyComparator(lastKey, currentKey) == comparison.Equal {
				continue
			}
		}

		firstIndex := lastIndex + 1
		lastIndex = i - 1

		validValues := mappedValues[firstIndex : lastIndex+1]

		reducerData := reducingDataGroup[KeyOut, ValueOut]{
			key:    lastKey,
			values: validValues,
		}

		readyDataPool <- reducerData
	}

	close(readyDataPool)
}
