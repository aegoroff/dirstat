package module

import (
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"io"
	"time"
)

// NewTimeMeasureCommand wraps real execute function to measure it's execution time
func NewTimeMeasureCommand(c Command) Command {
	return func(path string, fs afero.Fs, w io.Writer, modules ...Module) {
		start := time.Now()

		c(path, fs, w, modules...)

		elapsed := time.Since(start)

		color.Fprintf(w, "\n\n<gray>Read taken:\t%v</>\n", elapsed)
	}
}
