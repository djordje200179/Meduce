package meduce

import (
	"github.com/djordje200179/extendedlibrary/misc/functions"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
)

func mergeMappedData[KeyOut, ValueOut any](
	keyComparator functions.Comparator[KeyOut], valueComparator functions.Comparator[ValueOut],
	keysArrays [][]KeyOut, valuesArrays [][]ValueOut,
) ([]KeyOut, []ValueOut) {
	entriesCount := 0
	for _, keysArray := range keysArrays {
		entriesCount += len(keysArray)
	}

	mergedKeys := make([]KeyOut, 0, entriesCount)
	mergedValues := make([]ValueOut, 0, entriesCount)

	indices := make([]int, len(keysArrays))
	for i := 0; i < entriesCount; i++ {
		minIndex := -1
		var minKey KeyOut
		var minValue ValueOut

		for j, keys := range keysArrays {
			if indices[j] >= len(keys) {
				continue
			}

			currKey := keys[indices[j]]
			currValue := valuesArrays[j][indices[j]]
			if minIndex == -1 ||
				keyComparator(currKey, minKey) == comparison.FirstSmaller ||
				keyComparator(currKey, minKey) == comparison.Equal && valueComparator != nil && valueComparator(currValue, minValue) == comparison.FirstSmaller {
				minIndex = j
				minKey = currKey
				minValue = currValue
			}
		}

		mergedKeys = append(mergedKeys, minKey)
		mergedValues = append(mergedValues, minValue)

		indices[minIndex]++
	}

	return mergedKeys, mergedValues
}
