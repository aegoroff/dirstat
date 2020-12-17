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

type runner func(path string, fs afero.Fs, w io.Writer, modules ...module.Module)

func run(path string, c conf, modules ...module.Module) {
	var r runner
	{
		r = module.Execute
		r = newTimeMeasureR(r)
		r = newPrintMemoryR(r, *c.globals().showMemory)
		r = newPathCorrectionR(r)
	}

	r(path, c.fs(), c.w(), modules...)
}

func newTimeMeasureR(wrapped runner) runner {
	return func(path string, fs afero.Fs, w io.Writer, modules ...module.Module) {
		start := time.Now()

		wrapped(path, fs, w, modules...)

		elapsed := time.Since(start)

		color.Fprintf(w, "\n\n<gray>Read taken:\t%v</>\n", elapsed)
	}
}

func newPathCorrectionR(wrapped runner) runner {
	return func(path string, fs afero.Fs, w io.Writer, modules ...module.Module) {
		if len(path) == 0 {
			return
		}
		if (path)[len(path)-1] == ':' {
			path = filepath.Join(path, string(os.PathSeparator))
		}

		if _, err := fs.Stat(path); os.IsNotExist(err) {
			return
		}

		color.Fprintf(w, "Root: <red>%s</>\n\n", path)

		wrapped(path, fs, w, modules...)
	}
}

// newPrintMemoryR outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func newPrintMemoryR(wrapped runner, showMemory bool) runner {
	return func(path string, fs afero.Fs, w io.Writer, modules ...module.Module) {
		wrapped(path, fs, w, modules...)

		if !showMemory {
			return
		}
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		// For info on each, see: https://golang.org/pkg/runtime/#MemStats
		color.Fprintf(w, "\nAlloc = <gray>%s</>", humanize.IBytes(m.Alloc))
		color.Fprintf(w, "\tTotalAlloc = <gray>%s</>", humanize.IBytes(m.TotalAlloc))
		color.Fprintf(w, "\tSys = <gray>%s</>", humanize.IBytes(m.Sys))
		color.Fprintf(w, "\tNumGC = <gray>%v</>", m.NumGC)
		color.Fprintf(w, "\tNumGoRoutines = <gray>%v</>\n", runtime.NumGoroutine())
	}
}
