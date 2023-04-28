package sources

import (
	"github.com/djordje200179/extendedlibrary/misc"
	"github.com/djordje200179/meduce"
	"reflect"
	"runtime"
)

// AggregateDataSources aggregates multiple data sources into one.
func AggregateDataSources[K any, V any](dataSources ...meduce.Source[K, V]) meduce.Source[K, V] {
	cases := make([]reflect.SelectCase, len(dataSources))
	for i, dataSource := range dataSources {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(dataSource)}
	}
	remainingCases := len(cases)

	source := make(chan misc.Pair[K, V], runtime.NumCPU())
	go func() {
		for remainingCases > 0 {
			index, value, ok := reflect.Select(cases)
			if !ok {
				cases[index].Chan = reflect.ValueOf(nil)
				remainingCases--
				continue
			}

			source <- value.Interface().(misc.Pair[K, V])
		}

		close(source)
	}()

	return source
}
