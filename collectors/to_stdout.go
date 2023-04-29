package collectors

import (
	"fmt"
	"github.com/djordje200179/meduce"
	"sync"
)

type StdOutCollector[KeyOut, ValueOut any] struct {
	mutex sync.Mutex

	formatter Formatter[KeyOut, ValueOut]
}

// NewStdoutCollector creates a new FileCollector
// that writes key-value pairs to the standard output.
func NewStdoutCollector[KeyOut, ValueOut any]() meduce.Collector[KeyOut, ValueOut] {
	collector := &StdOutCollector[KeyOut, ValueOut]{}

	return collector
}

// NewStdoutCollectorWithFormatter creates a new FileCollector
// that writes key-value pairs to the standard output
// with the given formatter to format key-value pairs before writing them to a file.
func NewStdoutCollectorWithFormatter[KeyOut, ValueOut any](
	formatter Formatter[KeyOut, ValueOut],
) meduce.Collector[KeyOut, ValueOut] {
	collector := &StdOutCollector[KeyOut, ValueOut]{
		formatter: formatter,
	}

	return collector
}

func (collector *StdOutCollector[KeyOut, ValueOut]) Init() {
	collector.mutex.Lock()
}

func (collector *StdOutCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	var line string
	if collector.formatter != nil {
		line = collector.formatter(key, value)
	} else {
		line = fmt.Sprintf("%v: %v\n", key, value)
	}

	fmt.Print(line)
}

func (collector *StdOutCollector[KeyOut, ValueOut]) Finalize() {
	collector.mutex.Unlock()
}
