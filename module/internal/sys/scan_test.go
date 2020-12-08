package sys

import (
	"errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Scan(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/f/s", 0755)
	_ = afero.WriteFile(appFS, "/f/f.txt", []byte("123"), 0644)
	_ = afero.WriteFile(appFS, "/f/s/f.txt", []byte("1234"), 0644)
	var folders int
	var files int
	h := func(f *ScanEvent) {
		if f.File != nil {
			files++
		}
		if f.Folder != nil {
			folders++
		}
	}
	var handlers []ScanHandler
	handlers = append(handlers, h)

	// Act
	Scan("/", appFS, handlers)

	// Assert
	ass.Equal(2, files)
	ass.Equal(3, folders)
}

type errCloser struct{}

func (e *errCloser) Close() error {
	return errors.New("new error")
}

func Test_Close_ThatReturnsError(t *testing.T) {
	// Arrange
	ec := &errCloser{}

	// Act
	Close(ec)

	// Assert
}
