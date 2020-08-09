package module

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PathDecorator(t *testing.T) {
	var tests = []struct {
		root   string
		rr     bool
		result string
	}{
		{"/usr", true, "/local"},
		{"/usr", false, "/usr/local"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)

		pd := pathDecorator{
			removeRoot: test.rr,
			root:       test.root,
		}

		// Act
		decorated := pd.decorate(test.root + "/local")

		// Assert
		ass.Equal(test.result, decorated)
	}
}
