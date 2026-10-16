package module

import (
	"testing"

	"github.com/aegoroff/dirstat/scan"
	"github.com/stretchr/testify/assert"
)

func Test_DetailRanges_FileOnBoundaryInOneRange(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	all := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	details := newDetailsFile(newRanges(), all)
	h := newDetailFileHandler(details)

	// Act
	for _, sz := range []int64{100 * kbyte, mbyte, 10 * mbyte} {
		h.Handle(&scan.Event{File: &scan.FileEntry{Size: sz, Path: "/f"}})
	}

	// Assert
	ass.Len(details.enabled[1].files, 0)
	ass.Len(details.enabled[2].files, 1)
	ass.Len(details.enabled[3].files, 1)
	ass.Len(details.enabled[4].files, 1)
}
