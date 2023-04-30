package collectors

import (
	"fmt"
	"github.com/djordje200179/meduce"
	"sync"
)

var stdoutMutex sync.Mutex

type StdoutCollector[KeyOut, ValueOut any] struct {
	formatter Formatter[KeyOut, ValueOut]
}

// NewStdoutCollector creates a new FileCollector
// that writes key-value pairs to the standard output.
func NewStdoutCollector[KeyOut, ValueOut any]() meduce.Collector[KeyOut, ValueOut] {
	collector := StdoutCollector[KeyOut, ValueOut]{}

	return collector
}

// NewStdoutCollectorWithFormatter creates a new FileCollector
// that writes key-value pairs to the standard output
// with the given formatter to format key-value pairs before writing them to a file.
func NewStdoutCollectorWithFormatter[KeyOut, ValueOut any](
	formatter Formatter[KeyOut, ValueOut],
) meduce.Collector[KeyOut, ValueOut] {
	collector := StdoutCollector[KeyOut, ValueOut]{
		formatter: formatter,
	}

	return collector
}

func (collector StdoutCollector[KeyOut, ValueOut]) Init() {
	stdoutMutex.Lock()
}

func (collector StdoutCollector[KeyOut, ValueOut]) Collect(key KeyOut, value ValueOut) {
	var line string
	if collector.formatter != nil {
		line = collector.formatter(key, value)
	} else {
		line = fmt.Sprintf("%v: %v\n", key, value)
	}

	fmt.Print(line)
}

func (collector StdoutCollector[KeyOut, ValueOut]) Finalize() {
	stdoutMutex.Unlock()
}
