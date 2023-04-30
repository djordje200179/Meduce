package sources

import (
	"bufio"
	"github.com/djordje200179/extendedlibrary/misc"
	"github.com/djordje200179/meduce"
	"os"
)

// NewFileSource creates a new source that reads a file
// from the given path line by line.
func NewFileSource(path string) meduce.Source[int, string] {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	source := make(chan misc.Pair[int, string], 100)

	go func() {
		lineIndex := 0
		for scanner.Scan() {
			source <- misc.Pair[int, string]{lineIndex, scanner.Text()}
			lineIndex++
		}
		close(source)

		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	return source
}
