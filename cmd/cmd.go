package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
)

func beforeRunCmd(path *string, fs afero.Fs, w io.Writer) error {
	if _, err := fs.Stat(*path); os.IsNotExist(err) {
		return err
	}

	if (*path)[len(*path)-1] == ':' {
		*path = filepath.Join(*path, "\\")
	}

	color.Fprintf(w, "Root: <red>%s</>\n\n", *path)
	return nil
}
