package meduce

import (
	"fmt"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"strings"
	"sync"
)

type reducingDataGroup[KeyOut, ValueOut any] struct {
	key    KeyOut
	values []ValueOut
}

type reducingThread[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	*Process[KeyIn, ValueIn, KeyOut, ValueOut]

	reductionsCount  int
	collectionsCount int
}

func (thread *reducingThread[KeyIn, ValueIn, KeyOut, ValueOut]) run(
	dataPool <-chan reducingDataGroup[KeyOut, ValueOut],
	finishSignal *sync.WaitGroup,
) {
	for groupData := range dataPool {
		var reducedValue ValueOut
		if len(groupData.values) == 1 {
			reducedValue = groupData.values[0]
		} else {
			reducedValue = thread.Reducer(groupData.key, groupData.values)
		}

		thread.reductionsCount++

		if thread.Finalizer != nil {
			thread.Finalizer(groupData.key, &reducedValue)
		}

		if thread.Filter == nil || thread.Filter(groupData.key, &reducedValue) {
			thread.collect(groupData.key, reducedValue)
			thread.collectionsCount++
		}
	}

	if thread.Logger != nil {
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("Process %d: reducing thread finished\n", thread.uid))
		sb.WriteString(fmt.Sprintf("\t%d reductions finished\n", thread.reductionsCount))
		sb.WriteString(fmt.Sprintf("\t%d collections finished\n", thread.reductionsCount))

		thread.Logger.Print(sb.String())
	}

	finishSignal.Done()
}

func reducingDataGenerationThread[KeyOut, ValueOut any](
	keyComparator comparison.Comparator[KeyOut],
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
