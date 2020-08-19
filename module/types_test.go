package module

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/stretchr/testify/assert"
	"strconv"
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

func Test_heads(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	rs := newRanges()

	var tests = []struct {
		addNum bool
		match  string
	}{
		{true, "^\\s{0,1}\\d{1,2}\\..+"},
		{false, "^\\w.+"},
	}

	for _, test := range tests {
		// Act
		ss := rs.heads(test.addNum)

		// Assert
		for _, s := range ss {
			ass.Regexp(test.match, s)
		}
	}
}

func Test_insertToFixedTree(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	tree := newFixedTree(3)

	// Act
	for i := 0; i < 5; i++ {
		s := rbtree.NewString(strconv.Itoa(i))
		tree.insert(s)
	}

	// Assert
	ass.Equal(int64(3), tree.tree.Len())
	var r []string
	tree.tree.Descend(func(n rbtree.Node) bool {
		r = append(r, n.Key().String())
		return true
	})

	ass.ElementsMatch([]string{"4", "3", "2"}, r)
}