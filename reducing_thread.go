package meduce

import (
	"sync"
)

type reducingGroupData[KeyOut, ValueOut any] struct {
	key    KeyOut
	values []ValueOut
}

func reduceData[KeyOut, ValueOut any](
	reducer Reducer[KeyOut, ValueOut], finalizer Finalizer[KeyOut, ValueOut],
	dataPool chan reducingGroupData[KeyOut, ValueOut],
	write func(key KeyOut, value ValueOut), finishSignal *sync.WaitGroup,
) {
	for {
		var groupData reducingGroupData[KeyOut, ValueOut]
		groupData, ok := <-dataPool
		if !ok {
			break
		}

		reducedValue := reducer(groupData.key, groupData.values)
		if finalizer != nil {
			finalizer(groupData.key, &reducedValue)
		}

		write(groupData.key, reducedValue)
	}

	finishSignal.Done()
}
