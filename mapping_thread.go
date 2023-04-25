package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc/functions"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"sort"
	"sync"
)

func mappingThread[KeyIn, ValueIn, KeyOut, ValueOut any](
	keyComparator functions.Comparator[KeyOut],
	mapper Mapper[KeyIn, ValueIn, KeyOut, ValueOut], combiner Reducer[KeyOut, ValueOut],
	dataSource Source[KeyIn, ValueIn],
	keysPlace *[]KeyOut, valuesPlace *[]ValueOut,
	finishSignal *sync.WaitGroup,
) {
	mappedData := mappingThreadData[KeyOut, ValueOut]{
		keyComparator: keyComparator,
	}

	for entry := range dataSource {
		mapper(entry.First, entry.Second, mappedData.append)
	}

	sort.Sort(&mappedData)

	uniqueKeys, combinedValues := mappedData.combine(combiner)

	*keysPlace = uniqueKeys
	*valuesPlace = combinedValues

	finishSignal.Done()
}

type mappingThreadData[KeyOut, ValueOut any] struct {
	keyComparator functions.Comparator[KeyOut]

	keys   []KeyOut
	values []ValueOut
}

func (data *mappingThreadData[KeyOut, ValueOut]) append(key KeyOut, value ValueOut) {
	data.keys = append(data.keys, key)
	data.values = append(data.values, value)
}
func (data *mappingThreadData[KeyOut, ValueOut]) Len() int {
	return len(data.keys)
}

func (data *mappingThreadData[KeyOut, ValueOut]) Less(i, j int) bool {
	return data.keyComparator(data.keys[i], data.keys[j]) == comparison.FirstSmaller
}

func (data *mappingThreadData[KeyOut, ValueOut]) Swap(i, j int) {
	data.keys[i], data.keys[j] = data.keys[j], data.keys[i]
	data.values[i], data.values[j] = data.values[j], data.values[i]
}

func (data *mappingThreadData[KeyOut, ValueOut]) combine(reducer Reducer[KeyOut, ValueOut]) ([]KeyOut, []ValueOut) {
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

			if data.keyComparator(lastKey, currentKey) == comparison.Equal {
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
		reducedValue := reducer(lastKey, validValues)

		uniqueKeys = append(uniqueKeys, lastKey)
		combinedValues = append(combinedValues, reducedValue)
	}

	return uniqueKeys, combinedValues
}
