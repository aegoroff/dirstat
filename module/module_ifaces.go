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
	handler(evt *sys.ScanEvent)
}

type renderer interface {
	print(p printer)
	order() int
}
