//go:build windows || darwin

package scan

import (
	"testing"

	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_SkipTest(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/proc", 0755)
	_ = afero.WriteFile(fs, "/proc/f.txt", []byte("123"), 0644)

	th := testHandler{
		fipaths: make(c9s.HashSet[string]),
		fopaths: make(c9s.HashSet[string]),
		fp:      make([]string, 0),
	}

	// Act
	Scan("/", newFs(fs), &th)

	// Assert
	ass.Equal(1, th.files)
	ass.Equal(2, th.folders)
}
