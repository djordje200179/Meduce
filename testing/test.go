package main

import (
	"github.com/djordje200179/meduce"
	"github.com/djordje200179/meduce/collectors"
	"github.com/djordje200179/meduce/sources"
	"log"
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
		meduce.Config[int, string, string, int]{
			Mapper:  MapMovieToYear,
			Reducer: ReduceYearCounters,

			Source:    sources.NewFileSource("files/title_basics.tsv"),
			Collector: collectors.NewFileCollector[string, int]("output.txt"),

			Logger: log.Default(),
		},
	)

	process.Run()
}
