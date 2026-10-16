package module

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PathDecorator(t *testing.T) {
	const root = "/usr"

	var tests = []struct {
		name   string
		d      decorator
		result string
	}{
		{"remove root", &removeRootDecorator{root: root}, filepath.FromSlash("/local")},
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

func Test_RemoveRootDecorator_NotCleanRoot(t *testing.T) {
	var tests = []struct {
		root   string
		path   string
		result string
	}{
		{"./", filepath.Join("./", "a"), "/a"},
		{"./", filepath.Join("./", "sub"), "/sub"},
		{"./", "./", ""},
		{".", filepath.Join(".", ".hidden"), "/.hidden"},
		{"/usr/", filepath.Join("/usr/", "local"), "/local"},
		{"/usr/", "/usr/", ""},
		{"../x", filepath.Join("../x", "y"), "/y"},
	}

	for _, test := range tests {
		t.Run(test.root+" "+test.path, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			ctx := NewContext(10, true, test.root)

			// Act
			decorated := ctx.pd.decorate(test.path)

			// Assert
			ass.Equal(filepath.FromSlash(test.result), decorated)
		})
	}
}
