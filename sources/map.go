package sources

import (
	"github.com/djordje200179/extendedlibrary/misc"
	"github.com/djordje200179/meduce"
)

// NewMapSource creates a new source that iterates over
// the given map and returns its key-value pairs.
func NewMapSource[K comparable, V any](m map[K]V) meduce.Source[K, V] {
	source := make(chan misc.Pair[K, V], 100)

	go func() {
		for key, value := range m {
			source <- misc.Pair[K, V]{key, value}
		}
		close(source)
	}()

	return source
}
