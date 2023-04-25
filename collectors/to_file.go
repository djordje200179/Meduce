package collectors

import (
	"fmt"
	"github.com/djordje200179/meduce"
	"os"
)

type FileCollector[KeyOut, ValueOut any] struct {
	file *os.File
}

func NewFileCollector[KeyOut, ValueOut any](path string) meduce.Collector[KeyOut, ValueOut] {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	collector := FileCollector[KeyOut, ValueOut]{
		file: file,
	}

	return collector
}

func NewStdoutCollector[KeyOut, ValueOut any]() meduce.Collector[KeyOut, ValueOut] {
	collector := FileCollector[KeyOut, ValueOut]{
		file: os.Stdout,
	}

	return collector
}

func (collector FileCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	line := fmt.Sprintf("%v: %v\n", key, value)

	_, err := collector.file.WriteString(line)
	if err != nil {
		panic(err)
	}
}

func (collector FileCollector[KeyOut, ValueOut]) Finalize() {
	err := collector.file.Close()
	if err != nil {
		panic(err)
	}
}
