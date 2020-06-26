package cmd

import (
	"dirstat/module"
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"time"
)

func run(path string, c conf, modules ...module.Module) {
	timing := newTimeMeasureCommand(module.Execute)
	pathing := newPathCorrectionCommand(timing)
	pathing(path, c.fs(), c.w(), modules...)
}

func newTimeMeasureCommand(c module.Command) module.Command {
	return func(path string, fs afero.Fs, w io.Writer, modules ...module.Module) {
		start := time.Now()

		c(path, fs, w, modules...)

		elapsed := time.Since(start)

		color.Fprintf(w, "\n\n<gray>Read taken:\t%v</>\n", elapsed)
	}
}

func newPathCorrectionCommand(c module.Command) module.Command {
	return func(path string, fs afero.Fs, w io.Writer, modules ...module.Module) {
		if _, err := fs.Stat(path); os.IsNotExist(err) {
			return
		}

		if (path)[len(path)-1] == ':' {
			path = filepath.Join(path, "\\")
		}

		color.Fprintf(w, "Root: <red>%s</>\n\n", path)

		c(path, fs, w, modules...)
	}
}
