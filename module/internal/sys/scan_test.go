package sys

import (
	"errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testHandler struct {
	folders int
	files   int
}

func (t *testHandler) Handle(evt *ScanEvent) {
	if evt.File != nil {
		t.files++
	}
	if evt.Folder != nil {
		t.folders++
	}
}

func Test_Scan(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/f/s", 0755)
	_ = afero.WriteFile(appFS, "/f/f.txt", []byte("123"), 0644)
	_ = afero.WriteFile(appFS, "/f/s/f.txt", []byte("1234"), 0644)

	th := testHandler{}

	// Act
	Scan("/", appFS, &th)

	// Assert
	ass.Equal(2, th.files)
	ass.Equal(3, th.folders)
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
