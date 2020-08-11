package module

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRange_Contains(t *testing.T) {
	// Arrange
	r := Range{
		Min: 1,
		Max: 100,
	}
	var tests = []struct {
		val    int64
		result bool
	}{
		{1, true},
		{100, true},
		{50, true},
		{0, false},
		{101, false},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)

		// Act
		contains := r.Contains(test.val)

		// Assert
		ass.Equal(test.result, contains)
	}
}

func Test_countPercent(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	ti := totalInfo{
		FilesTotal: countSizeAggregate{Count: 100},
	}

	// Act
	r := ti.countPercent(10)

	// Assert
	ass.Equal(10.0, r)
}

func Test_sizePercent(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	ti := totalInfo{
		FilesTotal: countSizeAggregate{Size: 100},
	}

	// Act
	r := ti.sizePercent(10)

	// Assert
	ass.Equal(10.0, r)
}
