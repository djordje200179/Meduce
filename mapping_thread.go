package meduce

import (
	"fmt"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"sort"
	"strings"
	"sync"
)

func mappingThread[KeyIn, ValueIn, KeyOut, ValueOut any](
	process *Process[KeyIn, ValueIn, KeyOut, ValueOut],
	keysPlace *[]KeyOut, valuesPlace *[]ValueOut,
	finishSignal *sync.WaitGroup,
) {
	mappedData := mappingThreadData[KeyIn, ValueIn, KeyOut, ValueOut]{
		Process: process,
	}

	mappingsCount := 0

	for pair := range process.Source {
		process.Mapper(pair.First, pair.Second, mappedData.append)
		mappingsCount++
	}

	sort.Sort(&mappedData)

	uniqueKeys, combinedValues := mappedData.combine()

	*keysPlace = uniqueKeys
	*valuesPlace = combinedValues

	if process.Logger != nil {
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("Process %d: mapping thread finished\n", process.uid))
		sb.WriteString(fmt.Sprintf("\t%d mappings finished\n", mappingsCount))
		sb.WriteString(fmt.Sprintf("\t%d emmited key-value pairs\n", mappedData.Len()))
		sb.WriteString(fmt.Sprintf("\t%d unique keys\n", len(uniqueKeys)))

		process.Logger.Print(sb.String())
	}

	finishSignal.Done()
}

type mappingThreadData[KeyIn, ValueIn, KeyOut, ValueOut any] struct {
	*Process[KeyIn, ValueIn, KeyOut, ValueOut]

	keys   []KeyOut
	values []ValueOut
}

func (data *mappingThreadData[KeyIn, ValueIn, KeyOut, ValueOut]) append(key KeyOut, value ValueOut) {
	data.keys = append(data.keys, key)
	data.values = append(data.values, value)
}
func (data *mappingThreadData[KeyIn, ValueIn, KeyOut, ValueOut]) Len() int {
	return len(data.keys)
}

func (data *mappingThreadData[KeyIn, ValueIn, KeyOut, ValueOut]) Less(i, j int) bool {
	keyComparisonResult := data.KeyComparator(data.keys[i], data.keys[j])

	if keyComparisonResult == comparison.FirstSmaller {
		return true
	} else if keyComparisonResult == comparison.Equal {
		if data.ValueComparator == nil {
			return false
		}
		return data.ValueComparator(data.values[i], data.values[j]) == comparison.FirstSmaller
	} else {
		return false
	}
}

func (data *mappingThreadData[KeyIn, ValueIn, KeyOut, ValueOut]) Swap(i, j int) {
	data.keys[i], data.keys[j] = data.keys[j], data.keys[i]
	data.values[i], data.values[j] = data.values[j], data.values[i]
}

func (data *mappingThreadData[KeyIn, ValueIn, KeyOut, ValueOut]) combine() ([]KeyOut, []ValueOut) {
	if len(data.keys) == 0 {
		return nil, nil
	}

	uniqueKeys := make([]KeyOut, 0)
	combinedValues := make([]ValueOut, 0)

	lastIndex := -1
	for i := 1; i <= data.Len(); i++ {
		lastKey := data.keys[i-1]

		if i != data.Len() {
			currentKey := data.keys[i]

			if data.KeyComparator(lastKey, currentKey) == comparison.Equal {
				continue
			}
		}

		firstIndex := lastIndex + 1
		lastIndex = i - 1

		if firstIndex == lastIndex {
			value := data.values[firstIndex]
			uniqueKeys = append(uniqueKeys, lastKey)
			combinedValues = append(combinedValues, value)

			continue
		}

		validValues := data.values[firstIndex : lastIndex+1]
		reducedValue := data.Config.Reducer(lastKey, validValues)

		uniqueKeys = append(uniqueKeys, lastKey)
		combinedValues = append(combinedValues, reducedValue)
	}

	return uniqueKeys, combinedValues
}
