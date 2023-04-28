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
The paradigm is pretty simple. You only need two (or three) functions to process
all of your data:

1. `func Mapper(key KeyIn, value ValueIn, emit Emitter[KeyOut, ValueOut])`  
This function maps data that you supplied to key-value pairs. 
For each piece of data you can emit as many key-value pairs 
as you want by calling `emit` function. 

2. `func Reducer(key KeyOut, values []ValueOut) ValueOut`  
This function reduces all values that were mapped to the same key. 
It will be called many times to reduce local data, and then once more to reduce
all data from all threads. Because of that it should be idempotent and 
have no side effects.

3. `func Finalizer(key KeyOut, valueRef *ValueOut) ValueOut` (optional)  
This function is called after all data was processed and reduced. It is used to 
calculate final results in values that were reduced. It receives pointer to value,
so you should modify it in-place.  
If you don't need this function you can pass `nil` instead.

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
To start data processing you firstly need to create an `Process` object. That can be 
accomplished by using some constructor functions. The most general one is:
```go
func NewProcess[KeyIn, ValueIn, KeyOut, ValueOut any](
	keyComparator functions.Comparator[KeyOut],
	valueComparator functions.Comparator[ValueOut],
	
	mapper Mapper[KeyIn, ValueIn, KeyOut, ValueOut], 
	reducer Reducer[KeyOut, ValueOut], 
	finalizer Finalizer[KeyOut, ValueOut],
	
	dataSource Source[KeyIn, ValueIn], 
	collector Collector[KeyOut, ValueOut],
) *Process[KeyIn, ValueIn, KeyOut, ValueOut]
```

After creating the process you can start it either synchronously or asynchronously. And if you
start it asynchronously, you can wait for it to finish by calling `WaitToFinish()` method.
```go
go process.Run()
process.WaitToFinish()
```

## Example
In this example we are using IMDB title_basics dataset to find out how many movies
were released each year.

```go
package main

import (
	"github.com/djordje200179/meduce"
	"github.com/djordje200179/meduce/collectors"
	"github.com/djordje200179/meduce/sources"
	"strings"
)

func MapMovieToYear(_ int, line string, emit meduce.Emitter[string, int]) {
	values := strings.Split(line, "\t")
	year := values[5]

	emit(year, 1)
}

func ReduceYearCounters(_ string, counters []int) int {
	count := 0
	for _, value := range counters {
		count += value
	}

	return count
}

func main() {
	process := meduce.NewDefaultProcess(
		MapMovieToYear, ReduceYearCounters, nil,
		sources.NewFileSource("files/title_basics.tsv"),
		collectors.NewFileCollector[string, int]("output.txt"),
	)

	process.Run(true)
}
```
