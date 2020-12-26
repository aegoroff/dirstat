package module

import (
	"dirstat/module/internal/sys"
	"github.com/spf13/afero"
	"io"
)

// Module defines working modules interface
type Module interface {
	handlers() []sys.Handler
	renderers() []renderer
}

type renderer interface {
	print(p printer)
	order() int
}

// Execute runs modules over path specified
func Execute(path string, fs afero.Fs, w io.Writer, modules ...Module) {
	var renderers []renderer
	var handlers []sys.Handler

	for _, m := range modules {
		renderers = append(renderers, m.renderers()...)
		handlers = append(handlers, m.handlers()...)
	}

	sys.Scan(path, fs, handlers...)

	render(w, renderers)
}

// NewFoldersModule creates new folders module
func NewFoldersModule(ctx *Context, order int) Module {
	work := newFoldersHandler(ctx)
	rend := newFoldersRenderer(work, order)
	return newModule(rend, newOnlyFoldersWorker(work))
}

// NewTopFilesModule creates new top files statistic module
func NewTopFilesModule(ctx *Context, order int) Module {
	work := newTopFilesWorker(ctx.top, ctx.pd)
	rend := newTopFilesRenderer(work, order)

	m := newModule(rend, newOnlyFilesWorker(work))
	return m
}

// NewDetailFileModule creates new file statistic by file size range module
func NewDetailFileModule(ctx *Context, order int, enabledRanges []int) Module {
	// Do nothing if verbose not enabled
	if len(enabledRanges) == 0 {
		return newVoidModule()
	}
	work := newDetailFileHandler(newRanges(), enabledRanges, ctx.pd)
	rend := newDetailFileRenderer(work, order)

	m := newModule(rend, newOnlyFilesWorker(work))
	return m
}

// NewBenfordFileModule creates new file size bendford statistic
func NewBenfordFileModule(ctx *Context, order int) Module {
	work := newBenfordFileHandler(ctx)
	rend := newBenfordFileRenderer(work, order)

	m := newModule(rend, newOnlyFilesWorker(work))
	return m
}

// NewExtensionModule creates new file extensions statistic module
func NewExtensionModule(ctx *Context, order int) Module {
	rend := newExtRenderer(ctx, order)
	m := newModule(rend)
	return m
}

// NewAggregateFileModule creates new total file statistic module
func NewAggregateFileModule(ctx *Context, order int) Module {
	work := newAggregateFileHandler(newRanges())
	rend := newAggregateFileRenderer(ctx, work, order)

	m := newModule(rend, newOnlyFilesWorker(work))
	return m
}

// NewTotalModule creates new total statistic module
func NewTotalModule(ctx *Context, order int) Module {
	workFile := newTotalFileWorker(ctx)
	workFold := newTotalFolderWorker(ctx)
	rend := newTotalRenderer(ctx, order)

	m := newModule(rend, newOnlyFilesWorker(workFile), newOnlyFoldersWorker(workFold))
	return m
}

type module struct {
	wks []sys.Handler
	rnd []renderer
}

func (m *module) handlers() []sys.Handler {
	return m.wks
}

func (m *module) renderers() []renderer {
	return m.rnd
}

// newVoidModule creates module that do nothing
func newVoidModule() Module {
	return &module{
		[]sys.Handler{},
		[]renderer{},
	}
}

func newModule(r renderer, w ...sys.Handler) Module {
	m := module{
		[]sys.Handler{},
		[]renderer{r},
	}
	m.wks = append(m.wks, w...)
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

type baseRenderer struct {
	ord int
}

func newBaseRenderer(order int) *baseRenderer {
	return &baseRenderer{ord: order}
}

func (br *baseRenderer) order() int {
	return br.ord
}
