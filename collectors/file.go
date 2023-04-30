package collectors

import (
	"fmt"
	"github.com/djordje200179/meduce"
	"os"
)

// Formatter is a function that formats key and value into a string.
// It is used by FileCollector to format key and value before writing them to a file.
type Formatter[KeyOut, ValueOut any] func(key KeyOut, value ValueOut) string

// FileCollector is a collector that writes key-value pairs to a file.
type FileCollector[KeyOut, ValueOut any] struct {
	file      *os.File
	formatter Formatter[KeyOut, ValueOut]
}

// NewFileCollector creates a new FileCollector
// that writes key-value pairs to a file at the given path.
func NewFileCollector[KeyOut, ValueOut any](path string) FileCollector[KeyOut, ValueOut] {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	collector := FileCollector[KeyOut, ValueOut]{
		file: file,
	}

	return collector
}

// NewFileCollectorWithFormatter creates a new FileCollector
// that writes key-value pairs to a file at the given path
// with the given formatter to format key-value pairs before writing them to a file.
func NewFileCollectorWithFormatter[KeyOut, ValueOut any](
	path string,
	formatter Formatter[KeyOut, ValueOut],
) meduce.Collector[KeyOut, ValueOut] {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	collector := FileCollector[KeyOut, ValueOut]{
		file:      file,
		formatter: formatter,
	}

	return collector
}

func (collector FileCollector[KeyOut, ValueOut]) Init() {

}

func (collector FileCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	var line string
	if collector.formatter != nil {
		line = collector.formatter(key, value)
	} else {
		line = fmt.Sprintf("%v: %v\n", key, value)
	}

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
