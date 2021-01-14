package module

import (
	"fmt"
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
)

// Module defines working modules interface
type Module interface {
	handlers() []scan.Handler
	renderers() []renderer
}

// Folder interface
type folderI interface {
	fmt.Stringer
	Size() int64
	Count() int64
}

type renderer interface {
	ordered
	render(p out.Printer)
}

type ordered interface {
	order() int
}
