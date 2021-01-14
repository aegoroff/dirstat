package module

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PathDecorator(t *testing.T) {
	const root = "/usr"

	var tests = []struct {
		name   string
		d      decorator
		result string
	}{
		{"remove root", &removeRootDecorator{root: root}, "/local"},
		{"dont change", &noChangeDecorator{}, "/usr/local"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)

			// Act
			decorated := test.d.decorate(root + "/local")

			// Assert
			ass.Equal(test.result, decorated)
		})
	}
}
