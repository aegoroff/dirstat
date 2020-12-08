package module

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile_EqualTo(t *testing.T) {
	var tests = []struct {
		size1  int64
		size2  int64
		result bool
	}{
		{0, 0, true},
		{1, 1, true},
		{1, 0, false},
		{0, 1, false},
		{2, 1, false},
		{1, 2, false},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		f1 := &file{
			path: "/",
			size: test.size1,
			pd:   nil,
		}
		f2 := &file{
			path: "/f",
			size: test.size2,
			pd:   nil,
		}

		// Act
		result := f1.EqualTo(f2)

		// Assert
		ass.Equal(test.result, result)
	}
}

func TestFile_LessThan(t *testing.T) {
	var tests = []struct {
		size1  int64
		size2  int64
		result bool
	}{
		{0, 0, false},
		{1, 1, false},
		{1, 0, false},
		{0, 1, true},
		{2, 1, false},
		{1, 2, true},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		f1 := &file{
			path: "/",
			size: test.size1,
			pd:   nil,
		}
		f2 := &file{
			path: "/f",
			size: test.size2,
			pd:   nil,
		}

		// Act
		result := f1.LessThan(f2)

		// Assert
		ass.Equal(test.result, result)
	}
}

func TestFile_String(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	f := &file{
		path: "/",
		size: 0,
		pd:   nil,
	}

	// Act
	result := f.String()

	// Assert
	ass.Equal("/", result)
}

func TestFile_String_PathDecorating(t *testing.T) {
	// Arrange
	pd := removeRootDecorator{
		root: "/usr",
	}
	ass := assert.New(t)
	f := &file{
		path: "/usr/local",
		size: 0,
		pd:   &pd,
	}

	// Act
	result := f.String()

	// Assert
	ass.Equal("/local", result)
}
