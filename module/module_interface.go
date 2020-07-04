package module

import (
	"dirstat/module/internal/sys"
	"github.com/spf13/afero"
	"io"
)

// Module defines working modules interface
type Module interface {
	workers() []worker
	renderers() []renderer
}

// Runner defines executable function prototype
type Runner func(path string, fs afero.Fs, w io.Writer, modules ...Module)

type worker interface {
	initer
	finalizer
	handler(evt *sys.ScanEvent)
}

type initer interface {
	init()
}

type finalizer interface {
	finalize()
}

type renderer interface {
	print(p printer)
}
