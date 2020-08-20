package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RunAll(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/f/s", 0755)
	_ = afero.WriteFile(appFS, "/f/f.txt", []byte("123"), 0644)
	_ = afero.WriteFile(appFS, "/f/s/f.txt", []byte("1234"), 0644)
	w := bytes.NewBufferString("")

	// Act
	Execute(appFS, w, "a", "-p", "/")

	// Assert
	ass.Greater(w.Len(), 0)
}
