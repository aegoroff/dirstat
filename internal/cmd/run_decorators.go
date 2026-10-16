package cmd

import (
	"github.com/aegoroff/dirstat/internal/module"
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"runtime"
	"time"
)

type runner func(path string, fs afero.Fs, p out.Printer, modules ...module.Module)

func newTimeMeasureR(wrapped runner) runner {
	return func(path string, fs afero.Fs, p out.Printer, modules ...module.Module) {
		start := time.Now()

		wrapped(path, fs, p, modules...)

		elapsed := time.Since(start)

		p.Cprint("\n\n<gray>Read taken:\t%v</>\n", elapsed)
	}
}

func newPrintRootR(wrapped runner) runner {
	return func(path string, fs afero.Fs, p out.Printer, modules ...module.Module) {
		p.Cprint("Root: <red>%s</>\n\n", path)

		wrapped(path, fs, p, modules...)
	}
}

// newPrintMemoryR outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func newPrintMemoryR(wrapped runner) runner {
	return func(path string, fs afero.Fs, p out.Printer, modules ...module.Module) {
		wrapped(path, fs, p, modules...)

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		// For info on each, see: https://golang.org/pkg/runtime/#MemStats
		p.Cprint("\nAlloc = <gray>%s</>", humanize.IBytes(m.Alloc))
		p.Cprint("\tTotalAlloc = <gray>%s</>", humanize.IBytes(m.TotalAlloc))
		p.Cprint("\tSys = <gray>%s</>", humanize.IBytes(m.Sys))
		p.Cprint("\tNumGC = <gray>%v</>", m.NumGC)
		p.Cprint("\tNumGoRoutines = <gray>%v</>\n", runtime.NumGoroutine())
	}
}
