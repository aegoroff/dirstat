package module

import (
	"dirstat/module/internal/sys"
)

// Module defines working modules interface
type Module interface {
	workers() []worker
	renderers() []renderer
}

type worker interface {
	initer
	finalizer
	handlerer
}

type handlerer interface {
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
	order() int
}
