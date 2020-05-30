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
	addWorker(w worker)
	addRenderer(r renderer)
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

func (m *module) addWorker(w worker) {
	m.wks = append(m.wks, w)
}

func (m *module) addRenderer(r renderer) {
	m.rnd = append(m.rnd, r)
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

// NewContext creates new module's context that needed to create new modules
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
func Execute(path string, fs afero.Fs, w io.Writer, modules ...Module) {
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

// NewFoldersModule creates new folders module
func NewFoldersModule(ctx *Context) Module {
	work := newFoldersWorker(ctx)
	rend := newFoldersRenderer(work)

	m := newModule()
	m.addWorker(work)
	m.addRenderer(rend)
	return m
}

// NewFoldersHiddenModule creates new folders module
// that has disabled output
func NewFoldersHiddenModule(ctx *Context) Module {
	work := newFoldersWorker(ctx)
	m := newModule()
	m.addWorker(work)
	return m
}

// NewTopFilesModule creates new top files statistic module
func NewTopFilesModule(c *Context) Module {
	work := newTopFilesWorker(c.top)
	rend := newTopFilesRenderer(work)
	m := newModule()
	m.addWorker(work)
	m.addRenderer(rend)
	return m
}

// NewRangeModule creates new file statistic by file size range module
func NewRangeModule(ctx *Context, verbose bool, enabledRanges []int) Module {
	work := newRangeWorker(ctx, verbose, enabledRanges)
	rend := newRangeRenderer(work)
	m := newModule()
	m.addWorker(work)
	m.addRenderer(rend)
	return m
}

// NewExtensionModule creates new file extensions statistic module
func NewExtensionModule(ctx *Context) Module {
	work := newExtWorker(ctx)
	rend := newExtRenderer(work)
	m := newModule()
	m.addWorker(work)
	m.addRenderer(rend)
	return m
}

// NewExtensionHiddenModule creates new file extensions statistic module
// that has disabled output
func NewExtensionHiddenModule(ctx *Context) Module {
	work := newExtWorker(ctx)
	m := newModule()
	m.addWorker(work)
	return m
}

// NewTotalFileModule creates new total file statistic module
func NewTotalFileModule(ctx *Context) Module {
	r := newTotalFileRenderer(ctx)

	m := newModule()
	m.addRenderer(r)
	return m
}

// NewTotalModule creates new total statistic module
func NewTotalModule(ctx *Context) Module {
	work := newTotalWorker(ctx)
	rend := newTotalRenderer(work)

	m := newModule()
	m.addWorker(work)
	m.addRenderer(rend)
	return m
}

func newModule() Module {
	m := module{
		[]worker{},
		[]renderer{},
	}
	return &m
}
