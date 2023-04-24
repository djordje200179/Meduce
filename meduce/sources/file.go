package sources

import (
	"bufio"
	"github.com/djordje200179/extendedlibrary/misc"
	"github.com/djordje200179/meduce"
	"log"
	"os"
)

func NewFileSource(path string) meduce.Source[int, string] {
	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
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
	}()

	return source
}
