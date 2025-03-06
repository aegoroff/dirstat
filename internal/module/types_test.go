package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange_Floor(t *testing.T) {
	// Arrange
	tree := newRanges()

	var tests = []struct {
		name   string
		val    int64
		result int64
	}{
		{"-1", -1, 1},
		{"0", 0, 1},
		{"1", 1, 1},
		{"100*kbyte + 10", 100*kbyte + 10, 2},
		{"mbyte + 10", mbyte + 10, 3},
		{"10*mbyte + 10", 10*mbyte + 10, 4},
		{"100*mbyte + 10", 100*mbyte + 10, 5},
		{"gbyte + 10", gbyte + 10, 6},
		{"10*gbyte + 10", 10*gbyte + 10, 7},
		{"100*gbyte + 10", 100*gbyte + 10, 8},
		{"tbyte + 10", tbyte + 10, 9},
		{"10*tbyte + 10", 10*tbyte + 10, 10},
		{"pbyte + 10", pbyte + 10, 10},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			k := &Range{Min: test.val, Max: test.val}

			// Act
			r, ok := tree.Floor(k)

			// Assert
			ass.True(ok)
			expect, _ := tree.OrderStatisticSelect(test.result)
			ass.True(expect.Key().Equal(r))
		})
	}
}

func TestRange_Equal(t *testing.T) {
	// Arrange
	r := Range{
		Min: 0,
		Max: 100,
	}
	var tests = []struct {
		name   string
		val    int64
		result bool
	}{
		{"-1", -1, false},
		{"0", 0, true},
		{"1", 1, false},
		{"50", 50, false},
		{"100", 100, false},
		{"101", 101, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			k := &Range{Min: test.val, Max: test.val}

			// Act
			contains := r.Equal(k)

			// Assert
			ass.Equal(test.result, contains)
		})
	}
}

func TestRange_Less(t *testing.T) {
	// Arrange
	r := &Range{
		Min: 0,
		Max: 100,
	}
	var tests = []struct {
		name   string
		val    int64
		result bool
	}{
		{"-1", -1, true},
		{"0", 0, false},
		{"1", 1, false},
		{"50", 50, false},
		{"100", 100, false},
		{"101", 101, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			k := &Range{Min: test.val, Max: test.val}

			// Act
			contains := k.Less(r)

			// Assert
			ass.Equal(test.result, contains)
		})
	}
}

func TestRange_Contains(t *testing.T) {
	// Arrange
	r := Range{
		Min: 1,
		Max: 100,
	}
	var tests = []struct {
		name   string
		val    int64
		result bool
	}{
		{"-1", -1, false},
		{"1", 1, true},
		{"100", 100, true},
		{"50", 50, true},
		{"0", 0, false},
		{"101", 101, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)

			// Act
			contains := r.Contains(test.val)

			// Assert
			ass.Equal(test.result, contains)
		})
	}
}

func Test_countPercent(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	ti := totalInfo{
		FilesTotal: 100,
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
		FilesSize: 100,
	}

	// Act
	r := ti.sizePercent(10)

	// Assert
	ass.Equal(10.0, r)
}

func Test_percent(t *testing.T) {
	var tests = []struct {
		name     string
		val      float64
		total    float64
		expected float64
	}{
		{"0/100 = 0", 0, 100, 0},
		{"10/100 = 10", 10, 100, 10},
		{"10/0 = 0", 10, 0, 0},
		{"33.5/100 = 33.5", 33.5, 100, 33.5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)

			// Act
			result := percent(test.val, test.total)

			// Assert
			ass.Equal(test.expected, result)
		})
	}
}

func Test_TableWriter_sizeTransformer(t *testing.T) {
	var tests = []struct {
		val      any
		expected string
	}{
		{int64(10), "10 B"},
		{uint64(10), "10 B"},
		{"10", ""},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			tw := &tableWriter{}

			// Act
			result := tw.sizeTransformer(test.val)

			// Assert
			ass.Equal(test.expected, result)
		})
	}
}
