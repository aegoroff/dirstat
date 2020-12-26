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
	fc := newFolders(ctx.top)
	handler := newFoldersHandler(fc, ctx.pd)
	rend := newFoldersRenderer(fc, ctx, order)
	return newModule(rend, handler)
}

// NewTopFilesModule creates new top files statistic module
func NewTopFilesModule(ctx *Context, order int) Module {
	tf := newTopFiles(ctx.top)
	handler := newTopFilesHandler(tf, ctx.pd)
	rend := newTopFilesRenderer(tf, order)

	m := newModule(rend, handler)
	return m
}

// NewDetailFileModule creates new file statistic by file size range module
func NewDetailFileModule(ctx *Context, order int, enabledRanges []int) Module {
	// Do nothing if verbose not enabled
	if len(enabledRanges) == 0 {
		return newVoidModule()
	}
	details := newDetailsFile(newRanges())
	handler := newDetailFileHandler(details, enabledRanges, ctx.pd)
	rend := newDetailFileRenderer(details, order)

	m := newModule(rend, handler)
	return m
}

// NewBenfordFileModule creates new file size bendford statistic
func NewBenfordFileModule(ctx *Context, order int) Module {
	bf := newBenfordFile()
	handler := newBenfordFileHandler(bf)
	rend := newBenfordFileRenderer(ctx, bf, order)

	m := newModule(rend, handler)
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
	af := newAggregateFile(newRanges())
	handler := newAggregateFileHandler(af)
	rend := newAggregateFileRenderer(ctx, af, order)

	m := newModule(rend, handler)
	return m
}

// NewTotalModule creates new total statistic module
func NewTotalModule(ctx *Context, order int) Module {
	fi := newTotalFileHandler(ctx)
	fo := newTotalFolderHandler(ctx)
	rend := newTotalRenderer(ctx, order)

	m := newModule(rend, fi, fo)
	return m
}

type module struct {
	hlers []sys.Handler
	rnd   []renderer
}

func (m *module) handlers() []sys.Handler {
	return m.hlers
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

func newModule(r renderer, handlers ...sys.Handler) Module {
	m := module{
		[]sys.Handler{},
		[]renderer{r},
	}
	m.hlers = append(m.hlers, handlers...)
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
