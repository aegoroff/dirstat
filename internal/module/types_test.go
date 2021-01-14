package module

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRange_Search(t *testing.T) {
	// Arrange
	tree := newRanges()

	var tests = []struct {
		name   string
		val    int64
		result bool
	}{
		{"-1", -1, false},
		{"1", 1, true},
		{"100*kbyte + 10", 100*kbyte + 10, true},
		{"mbyte + 10", mbyte + 10, true},
		{"10*mbyte + 10", 10*mbyte + 10, true},
		{"100*mbyte + 10", 100*mbyte + 10, true},
		{"gbyte + 10", gbyte + 10, true},
		{"10*gbyte + 10", 10*gbyte + 10, true},
		{"100*gbyte + 10", 100*gbyte + 10, true},
		{"tbyte + 10", tbyte + 10, true},
		{"10*tbyte + 10", 10*tbyte + 10, true},
		{"pbyte + 10", pbyte + 10, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			k := &Range{Min: test.val, Max: test.val}

			// Act
			_, ok := tree.Search(k)

			// Assert
			ass.Equal(test.result, ok)
		})
	}
}

func TestRange_MatchingPredefinedRanges(t *testing.T) {
	// Arrange
	var tests = []struct {
		name string
		val  int64
		r    *Range
	}{
		{"0", 0, NewRange(0, 100*kbyte)},
		{"1", 1, NewRange(0, 100*kbyte)},
		{"100*kbyte", 100 * kbyte, NewRange(0, 100*kbyte)},
		{"100*kbyte + 1", 100*kbyte + 1, NewRange(100*kbyte+1, mbyte)},
		{"100*kbyte + 10", 100*kbyte + 10, NewRange(100*kbyte+1, mbyte)},
		{"mbyte", mbyte, NewRange(100*kbyte+1, mbyte)},
		{"mbyte + 1", mbyte + 1, NewRange(mbyte+1, 10*mbyte)},
		{"mbyte + 10", mbyte + 10, NewRange(mbyte+1, 10*mbyte)},
		{"10*mbyte", 10 * mbyte, NewRange(mbyte+1, 10*mbyte)},
		{"10*mbyte + 10", 10*mbyte + 10, NewRange(10*mbyte+1, 100*mbyte)},
		{"100*mbyte + 10", 100*mbyte + 10, NewRange(100*mbyte+1, gbyte)},
		{"gbyte + 10", gbyte + 10, NewRange(gbyte+1, 10*gbyte)},
		{"10*gbyte + 10", 10*gbyte + 10, NewRange(10*gbyte+1, 100*gbyte)},
		{"100*gbyte + 10", 100*gbyte + 10, NewRange(100*gbyte+1, tbyte)},
		{"tbyte + 10", tbyte + 10, NewRange(tbyte+1, 10*tbyte)},
		{"10*tbyte + 10", 10*tbyte + 10, NewRange(10*tbyte+1, pbyte)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			k := &Range{Min: test.val, Max: test.val}

			// Act
			less := test.r.Less(k)
			eq := test.r.Equal(k)

			// Assert
			ass.False(less)
			ass.True(eq)
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
		{"1", 1, true},
		{"100", 100, true},
		{"50", 50, true},
		{"0", 0, true},
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
		{"1", 1, false},
		{"100", 100, false},
		{"50", 50, false},
		{"0", 0, false},
		{"101", 101, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			k := &Range{Min: test.val, Max: test.val}

			// Act
			contains := r.Less(k)

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
		val      interface{}
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
