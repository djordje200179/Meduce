package sources

import (
	"github.com/djordje200179/extendedlibrary/misc"
	"github.com/djordje200179/meduce"
)

// NewSliceSource creates a new source that reads a slice
// and returns its elements with their indexes.
func NewSliceSource[T any](slice []T) meduce.Source[int, T] {
	source := make(chan misc.Pair[int, T], 100)

	go func() {
		for index, element := range slice {
			source <- misc.Pair[int, T]{index, element}
		}
		close(source)
	}()

	return source
}
