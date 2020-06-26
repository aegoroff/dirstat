package cmd

import (
	"dirstat/module"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func run(path string, c conf, modules ...module.Module) {
	timing := newTimeMeasureR(module.Execute)
	memory := newPrintMemoryR(timing)
	pathing := newPathCorrectionR(memory)
	pathing(path, c.fs(), c.w(), modules...)
}

func newTimeMeasureR(c module.Runner) module.Runner {
	return func(path string, fs afero.Fs, w io.Writer, modules ...module.Module) {
		start := time.Now()

		c(path, fs, w, modules...)

		elapsed := time.Since(start)

		color.Fprintf(w, "\n\n<gray>Read taken:\t%v</>\n", elapsed)
	}
}

func newPathCorrectionR(c module.Runner) module.Runner {
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

// newPrintMemoryR outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func newPrintMemoryR(c module.Runner) module.Runner {
	return func(path string, fs afero.Fs, w io.Writer, modules ...module.Module) {
		c(path, fs, w, modules...)

		if !showMemory {
			return
		}
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		// For info on each, see: https://golang.org/pkg/runtime/#MemStats
		color.Fprintf(w, "\nAlloc = <gray>%s</>", humanize.IBytes(m.Alloc))
		color.Fprintf(w, "\tTotalAlloc = <gray>%s</>", humanize.IBytes(m.TotalAlloc))
		color.Fprintf(w, "\tSys = <gray>%s</>", humanize.IBytes(m.Sys))
		color.Fprintf(w, "\tNumGC = <gray>%v</>\n", m.NumGC)
	}
}
