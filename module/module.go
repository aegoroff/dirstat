package module

import (
	"dirstat/module/internal/sys"
	"github.com/spf13/afero"
	"io"
)

// Execute runs modules over path specified
func Execute(path string, fs afero.Fs, w io.Writer, modules ...Module) {
	var renderers []renderer
	var workers []worker

	for _, m := range modules {
		renderers = append(renderers, m.renderers()...)
		workers = append(workers, m.workers()...)
	}

	handlers := make([]sys.ScanHandler, len(workers))
	for i, wo := range workers {
		wo.init()
		handlers[i] = wo.handler
	}

	sys.Scan(path, fs, handlers)

	for _, wo := range workers {
		wo.finalize()
	}

	render(w, renderers)
}

// NewFoldersModule creates new folders module
func NewFoldersModule(ctx *Context, order int) Module {
	work := newFoldersWorker(ctx)
	rend := newFoldersRenderer(work, order)
	return newModule(rend, work)
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
		return NewVoidModule()
	}
	work := newDetailFileWorker(newRanges(), enabledRanges, ctx.pd)
	rend := newDetailFileRenderer(work, order)

	m := newModule(rend, newOnlyFilesWorker(work))
	return m
}

// NewBenfordFileModule creates new file size bendford statistic
func NewBenfordFileModule(ctx *Context, order int) Module {
	work := newBenfordFileWorker(ctx)
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
	work := newAggregateFileWorker(newRanges())
	rend := newAggregateFileRenderer(ctx, work, order)

	m := newModule(rend, newOnlyFilesWorker(work))
	return m
}

// NewTotalModule creates new total statistic module
func NewTotalModule(ctx *Context, order int) Module {
	work := newTotalWorker(ctx)
	rend := newTotalRenderer(ctx, order)

	m := newModule(rend, work)
	return m
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

// NewVoidModule creates module that do nothing
func NewVoidModule() Module {
	return &module{
		[]worker{},
		[]renderer{},
	}
}

func newModule(r renderer, w ...worker) Module {
	m := module{
		[]worker{},
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

type voidInit struct{}

func (*voidInit) init() {}

type voidFinalize struct{}

func (*voidFinalize) finalize() {}

type baseRenderer struct {
	ord int
}

func newBaseRenderer(order int) *baseRenderer {
	return &baseRenderer{ord: order}
}

func (br *baseRenderer) order() int {
	return br.ord
}
