package meduce

import (
	"fmt"
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
		process.Logger.Printf("Process %d: %d mapping threads were started\n", process.uid, threadsCount)
	}

	allMappersFinished.Wait()

	if process.Logger != nil {
		process.Logger.Printf("Process %d: all mapping threads finished\n", process.uid)
	}

	var keysArray [][]KeyOut
	var valuesArray [][]ValueOut

	for i := range process.mappingThreads {
		keysArray = append(keysArray, process.mappingThreads[i].keys)
		valuesArray = append(valuesArray, process.mappingThreads[i].values)
	}

	process.mappedKeys, process.mappedValues = mergeMappingThreadsData(
		process.KeyComparator, process.ValueComparator,
		keysArray, valuesArray,
	)

	if process.Logger != nil {
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("Process %d: mapped data merged\n", process.uid))
		sb.WriteString(fmt.Sprintf("\t%d key-value pairs left\n", len(process.mappedKeys)))

		process.Logger.Print(sb.String())
	}
}
