package module

import (
	"testing"

	"github.com/aegoroff/dirstat/scan"
	"github.com/stretchr/testify/assert"
)

func Test_BenfordFirstDigit_LargeSizes(t *testing.T) {
	var tests = []struct {
		size  int64
		digit int
	}{
		{0, 0},
		{7, 7},
		{29886425451, 2},
		{100_000_000_000, 1},
		{200_000_000_000, 2},
		{2_000_000_000_000, 2},
		{999_999_999_999, 9},
		{1 << 50, 1},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		bf := newBenfordFile()
		h := newBenfordFileHandler(bf)

		// Act
		h.Handle(&scan.Event{File: &scan.FileEntry{Size: test.size}})

		// Assert
		ass.Equal(int64(1), bf.distribution[test.digit], "size %d", test.size)
	}
}
