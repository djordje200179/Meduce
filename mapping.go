package meduce

import (
	"fmt"
	"github.com/djordje200179/extendedlibrary/misc/functions/comparison"
	"runtime"
	"strings"
	"sync"
)

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) mapData() {
	threadsCount := runtime.NumCPU()

	var allMappersFinished sync.WaitGroup
	allMappersFinished.Add(threadsCount)

	process.mappingThreads = make([]mappingThread[KeyIn, ValueIn, KeyOut, ValueOut], threadsCount)

	for i := range process.mappingThreads {
		process.mappingThreads[i].Process = process

		go process.mappingThreads[i].run(&allMappersFinished)
	}

	if process.Logger != nil {
		var message string
		if threadsCount == 1 {
			message = "Process %d: 1 mapping thread was started\n"
		} else {
			message = "Process %d: %d mapping threads were started\n"
		}

		process.Logger.Printf(message, process.uid, threadsCount)
	}

	allMappersFinished.Wait()

	if process.Logger != nil {
		process.Logger.Printf("Process %d: all mapping threads finished\n", process.uid)
	}

	process.mergeMappedData()
	for i := range process.mappingThreads {
		process.mappingThreads[i].keys = nil
		process.mappingThreads[i].values = nil
	}

	if process.Logger != nil {
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("Process %d: mapped data merged\n", process.uid))
		sb.WriteString(fmt.Sprintf("\t%d key-value pairs left\n", len(process.mappedKeys)))

		process.Logger.Print(sb.String())
	}
}

func (process *Process[KeyIn, ValueIn, KeyOut, ValueOut]) mergeMappedData() {
	entriesCount := 0
	for _, thread := range process.mappingThreads {
		entriesCount += len(thread.keys)
	}

	process.mappedKeys = make([]KeyOut, entriesCount)
	process.mappedValues = make([]ValueOut, entriesCount)

	indices := make([]int, len(process.mappingThreads))
	for i := 0; i < entriesCount; i++ {
		minIndex := -1
		var minKey KeyOut
		var minValue ValueOut

		for j, thread := range process.mappingThreads {
			if indices[j] >= len(thread.keys) {
				continue
			}

			currKey := thread.keys[indices[j]]
			currValue := thread.values[indices[j]]
			if minIndex == -1 ||
				process.KeyComparator(currKey, minKey) == comparison.FirstSmaller ||
				process.KeyComparator(currKey, minKey) == comparison.Equal && process.ValueComparator != nil && process.ValueComparator(currValue, minValue) == comparison.FirstSmaller {
				minIndex = j
				minKey = currKey
				minValue = currValue
			}
		}

		process.mappedKeys[i] = minKey
		process.mappedValues[i] = minValue

		indices[minIndex]++
	}
}
