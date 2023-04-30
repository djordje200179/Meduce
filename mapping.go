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

	keysArrays := make([][]KeyOut, threadsCount)
	valuesArrays := make([][]ValueOut, threadsCount)

	for i := 0; i < threadsCount; i++ {
		go mappingThread(
			process,
			&keysArrays[i], &valuesArrays[i],
			&allMappersFinished,
		)
	}

	if process.Logger != nil {
		process.Logger.Printf("Process %d: %d mapping threads were started\n", process.uid, threadsCount)
	}

	allMappersFinished.Wait()

	if process.Logger != nil {
		process.Logger.Printf("Process %d: all mapping threads finished\n", process.uid)
	}

	process.mappedKeys, process.mappedValues = mergeMappedData(
		process.KeyComparator, process.ValueComparator,
		keysArrays, valuesArrays,
	)

	if process.Logger != nil {
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("Process %d: mapped data merged\n", process.uid))
		sb.WriteString(fmt.Sprintf("\t%d key-value pairs left\n", len(process.mappedKeys)))

		process.Logger.Print(sb.String())
	}
}
