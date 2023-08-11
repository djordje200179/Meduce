# Meduce

A library for processing huge amounts of data on your device by using 
MapReduce paradigm. It was inspired by Hadoop MapReduce and MongoDB MapReduce.

The goal of the library is to fully utilize all of your CPU cores
and maximize concurrent processing of data. It is not meant to be
used for distributed processing, but it can be used to process data
in parallel on a single machine.

It is written in Go 1.20, and it fully utilizes generic mechanics, 
so you don't need to worry about casts from `interface{}`.

## Usage

### Functions
The paradigm is pretty simple. You only need two (or four) functions to process
all of your data:

1. `func Mapper(key KeyIn, value ValueIn, emit Emitter[KeyOut, ValueOut])`
2. `func Reducer(key KeyOut, values []ValueOut) ValueOut`
3. `func Finalizer(key KeyOut, valueRef *ValueOut) ValueOut` _(optional)_
4. `func Filter(key KeyOut, valueRef *ValueOut) bool` _(optional)_

### Sources
Data is gathered from a channel named `Source`. You can use any channel, but most
commonly used ones are already predefined for you. And you can instantiate them
by calling suitable constructor functions:
1.	`func NewFileSource(path string) meduce.Source[int, string]`
2.  `func NewMapSource(m map[K]V) meduce.Source[K, V]`
3.  `func NewSliceSource(slice []T) meduce.Source[int, T]`

And if you need to read data from multiple sources, you can aggregate all of them
by calling a function that joins all sources into one:   
`func AggregateDataSources(dataSources ...meduce.Source[K, V]) meduce.Source[K, V]`

### Collectors
After all data is processed, it is collected by using a `Collector`. You can either
use predefined collectors (`FileCollector`, `MapCollector`, `ChannelCollector`) or
create your own that implements `Collector[K, V]` interface.

### Process
To start data processing, you firstly need to create an `Process` object. 
That can be accomplished by calling a constructor function.

After creating the process, you can start it either synchronously or asynchronously. And if you
start it asynchronously, you can wait for it to finish by calling `WaitToFinish()` method.
```go
go process.Run()
process.WaitToFinish()
```

### Links
You can link multiple processes together to create a pipeline.
Interconnected processes will share data between each other internally.

They will be executed in parallel, and you don't need to worry about
starting them manually. 
You only need to start the first one asynchronously, and the rest 
will be started automatically.

If you want to wait for all processes to finish, you can wait for the
last one to finish by calling `WaitToFinish()` method on it.

### Common reducers
In the `reducers` package, you can find some common reducers that 
you can use in your processes.

## Example
In this example, we're using IMDB title_basics dataset (can be found [here](https://datasets.imdbws.com/)) 
to find out in which year the most movies were released. 

We will create two linked processes. 
The first one will count how many movies were released in each year,
and the second one will find the year with the most movies.

```go
package main

import (
	"github.com/djordje200179/meduce"
	"github.com/djordje200179/meduce/collectors"
	"github.com/djordje200179/meduce/sources"
	"log"
	"strconv"
	"strings"
)

func MapMovieToYear(_ int, line string, emit meduce.Emitter[int, int]) {
	values := strings.Split(line, "\t")

	year, err := strconv.Atoi(values[5])
	if err != nil {
		return
	}

	emit(year, 1)
}

func ReduceYearCounters(_ int, counters []int) int {
	count := 0
	for _, value := range counters {
		count += value
	}

	return count
}

type YearInfo struct {
	Year  int
	Count int
}

func MapYearToInfo(year int, count int, emit meduce.Emitter[string, YearInfo]) {
	emit("max", YearInfo{year, count})
}

func ReduceMaxYear(_ string, infoList []YearInfo) YearInfo {
	var max YearInfo

	for _, info := range infoList {
		if max.Year == 0 || info.Count > max.Count {
			max = info
		}
	}

	return max
}

func main() {
	process1 := meduce.NewDefaultProcess(
		meduce.Config[int, string, int, int]{
			Mapper:  MapMovieToYear,
			Reducer: ReduceYearCounters,

			Source: sources.NewFileSource("files/title_basics.tsv"),

			Logger: log.Default(),
		},
	)

	maxValueCollector := collectors.NewSingleValueCollector[string, YearInfo]()

	process2 := meduce.NewDefaultProcess(
		meduce.Config[int, int, string, YearInfo]{
			Mapper:  MapYearToInfo,
			Reducer: reducers.NewMaxOrderedField[string, YearInfo, int](func(info YearInfo) int { return info.Count }),

			Collector: maxValueCollector,

			Logger: log.Default(),
		},
	)

	meduce.Link(process1, process2)

	go process1.Run()
	process2.WaitToFinish()

	maxValue := maxValueCollector.Value()
	fmt.Printf("Most movies (%d) were made in %d. year.\n", maxValue.Count, maxValue.Year)
}

```
