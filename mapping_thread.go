package meduce

import (
	"fmt"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"sort"
	"strings"
	"sync"
)

type mappingThread[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	*Process[KeyIn, ValueIn, KeyOut, ValueOut]

	keys   []KeyOut
	values []ValueOut

	mappingsCount     int
	emitsCount        int
	combinationsCount int
}

func (thread *mappingThread[KeyIn, ValueIn, KeyOut, ValueOut]) run(finishSignal *sync.WaitGroup) {
	for pair := range thread.Source {
		thread.Mapper(pair.First, pair.Second, thread.append)
		thread.mappingsCount++
	}

	thread.emitsCount = thread.Len()

	sort.Sort(thread)

	thread.combine()

	thread.combinationsCount = thread.Len()

	if thread.Logger != nil {
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("Process %d: mapping thread finished\n", thread.uid))
		sb.WriteString(fmt.Sprintf("\t%d mappings finished\n", thread.mappingsCount))
		sb.WriteString(fmt.Sprintf("\t%d emmited key-value pairs\n", thread.emitsCount))
		sb.WriteString(fmt.Sprintf("\t%d unique keys\n", thread.combinationsCount))

		thread.Logger.Print(sb.String())
	}

	finishSignal.Done()
}

func (thread *mappingThread[KeyIn, ValueIn, KeyOut, ValueOut]) append(key KeyOut, value ValueOut) {
	thread.keys = append(thread.keys, key)
	thread.values = append(thread.values, value)
}
func (thread *mappingThread[KeyIn, ValueIn, KeyOut, ValueOut]) Len() int {
	return len(thread.keys)
}

func (thread *mappingThread[KeyIn, ValueIn, KeyOut, ValueOut]) Less(i, j int) bool {
	keyComparisonResult := thread.KeyComparator(thread.keys[i], thread.keys[j])

	if keyComparisonResult == comparison.FirstSmaller {
		return true
	} else if keyComparisonResult == comparison.Equal {
		if thread.ValueComparator == nil {
			return false
		}
		return thread.ValueComparator(thread.values[i], thread.values[j]) == comparison.FirstSmaller
	} else {
		return false
	}
}

func (thread *mappingThread[KeyIn, ValueIn, KeyOut, ValueOut]) Swap(i, j int) {
	thread.keys[i], thread.keys[j] = thread.keys[j], thread.keys[i]
	thread.values[i], thread.values[j] = thread.values[j], thread.values[i]
}

func (thread *mappingThread[KeyIn, ValueIn, KeyOut, ValueOut]) combine() {
	if len(thread.keys) == 0 {
		return
	}

	uniqueKeys := make([]KeyOut, 0)
	combinedValues := make([]ValueOut, 0)

	lastIndex := -1
	for i := 1; i <= thread.Len(); i++ {
		lastKey := thread.keys[i-1]

		if i != thread.Len() {
			currentKey := thread.keys[i]

			if thread.KeyComparator(lastKey, currentKey) == comparison.Equal {
				continue
			}
		}

		firstIndex := lastIndex + 1
		lastIndex = i - 1

		if firstIndex == lastIndex {
			value := thread.values[firstIndex]
			uniqueKeys = append(uniqueKeys, lastKey)
			combinedValues = append(combinedValues, value)

			continue
		}

		validValues := thread.values[firstIndex : lastIndex+1]
		reducedValue := thread.Config.Reducer(lastKey, validValues)

		uniqueKeys = append(uniqueKeys, lastKey)
		combinedValues = append(combinedValues, reducedValue)
	}

	thread.keys = uniqueKeys
	thread.values = combinedValues
}
