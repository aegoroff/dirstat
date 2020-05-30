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

type module struct {
	wks []worker
	rnd []renderer
}

func (m *module) workers() []worker {
	return m.wks
}

func (m *module) renderers() []renderer {
	return m.rnd
}

type worker interface {
	init()
	handler(evt *sys.ScanEvent)
	finalize()
}

type renderer interface {
	print(p printer)
}

// Context defines modules context
type Context struct {
	total          *totalInfo
	rangeAggregate map[Range]fileStat
	top            int
}

// NewContext create new module's context that needed to create new modules
func NewContext(top int) *Context {
	total := totalInfo{}

	ctx := Context{
		total:          &total,
		rangeAggregate: make(map[Range]fileStat),
		top:            top,
	}
	return &ctx
}

// Execute runs modules over path specified
func Execute(path string, fs afero.Fs, w io.Writer, modules []Module) {
	var renderers []renderer
	var workers []worker

	for _, m := range modules {
		renderers = append(renderers, m.renderers()...)
		workers = append(workers, m.workers()...)
	}

	var handlers []sys.ScanHandler
	for _, wo := range workers {
		wo.init()
		handlers = append(handlers, wo.handler)
	}

	sys.Scan(path, fs, handlers)

	for _, m := range workers {
		m.finalize()
	}

	render(w, renderers)
}
