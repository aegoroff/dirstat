package module

import (
	"dirstat/module/internal/sys"
	"github.com/spf13/afero"
	"io"
)

// Context defines modules context
type Context struct {
	total *totalInfo
	top   int
	pd    *pathDecorator
}

// NewContext creates new module's context that needed to create new modules
func NewContext(top int, rr bool, root string) *Context {
	total := totalInfo{}

	ctx := Context{
		total: &total,
		top:   top,
		pd: &pathDecorator{
			removeRoot: rr,
			root:       root,
		},
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

	for _, wo := range workers {
		wo.finalize()
	}

	render(w, renderers)
}

// NewFoldersModule creates new folders module
func NewFoldersModule(ctx *Context, hideOutput bool) Module {
	work := newFoldersWorker(ctx)
	if hideOutput {
		return newModule(work)
	}
	rend := newFoldersRenderer(work)
	return newModule(work, rend)
}

// NewTopFilesModule creates new top files statistic module
func NewTopFilesModule(ctx *Context) Module {
	work := newTopFilesWorker(ctx.top, ctx.pd)
	rend := newTopFilesRenderer(work)
	m := newModule(work, rend)
	return m
}

// NewDetailFileModule creates new file statistic by file size range module
func NewDetailFileModule(ctx *Context, enabledRanges []int) Module {
	// Do nothing if verbose not enabled
	if len(enabledRanges) == 0 {
		return &module{
			[]worker{},
			[]renderer{},
		}
	}
	work := newDetailFileWorker(newRanges(), enabledRanges, ctx.pd)
	rend := newDetailFileRenderer(work)
	m := newModule(work, rend)
	return m
}

// NewBendfordFileModule creates new file size bendford statistic
func NewBendfordFileModule(ctx *Context, enabled bool) Module {
	if !enabled {
		return &module{
			[]worker{},
			[]renderer{},
		}
	}
	work := newBendfordFileWorker(ctx)
	rend := newBendfordFileRenderer(work)
	m := newModule(work, rend)
	return m
}

// NewExtensionModule creates new file extensions statistic module
func NewExtensionModule(ctx *Context, hideOutput bool) Module {
	work := newExtWorker(ctx)
	if hideOutput {
		return newModule(work)
	}
	rend := newExtRenderer(work, ctx.top)
	return newModule(work, rend)
}

// NewAggregateFileModule creates new total file statistic module
func NewAggregateFileModule(ctx *Context) Module {
	work := newAggregateFileWorker(newRanges())
	rend := newAggregateFileRenderer(ctx, work)

	m := newModule(work, rend)
	return m
}

// NewTotalModule creates new total statistic module
func NewTotalModule(ctx *Context) Module {
	rend := newTotalRenderer(ctx)

	m := module{
		[]worker{},
		[]renderer{rend},
	}
	return &m
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

type voidInit struct{}

func (*voidInit) init() {}

type voidFinalize struct{}

func (*voidFinalize) finalize() {}

func newModule(w worker, r ...renderer) Module {
	m := module{
		[]worker{w},
		[]renderer{},
	}
	m.rnd = append(m.rnd, r...)
	return &m
}

func newRanges() ranges {
	rs := []Range{
		{Min: 0, Max: 100 * kbyte},
		{Min: 100 * kbyte, Max: mbyte},
		{Min: mbyte, Max: 10 * mbyte},
		{Min: 10 * mbyte, Max: 100 * mbyte},
		{Min: 100 * mbyte, Max: gbyte},
		{Min: gbyte, Max: 10 * gbyte},
		{Min: 10 * gbyte, Max: 100 * gbyte},
		{Min: 100 * gbyte, Max: tbyte},
		{Min: tbyte, Max: 10 * tbyte},
		{Min: 10 * tbyte, Max: pbyte},
	}
	return rs
}
