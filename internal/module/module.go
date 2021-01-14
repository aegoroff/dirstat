package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/spf13/afero"
	"sort"
)

type filesystem struct {
	fs afero.Fs
}

func newFs(fs afero.Fs) scan.Filesystem {
	return &filesystem{fs: fs}
}

func (f *filesystem) Open(path string) (scan.File, error) {
	return f.fs.Open(path)
}

// Execute runs modules over path specified
func Execute(path string, fs afero.Fs, p out.Printer, modules ...Module) {
	var renderers []renderer
	var handlers []scan.Handler

	for _, m := range modules {
		renderers = append(renderers, m.renderers()...)
		handlers = append(handlers, m.handlers()...)
	}

	scan.Scan(path, newFs(fs), handlers...)

	render(renderers, p)
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
	handler := newTopFilesHandler(tf)
	rend := newTopFilesRenderer(tf, ctx.pd, order)

	m := newModule(rend, handler)
	return m
}

// NewDetailFileModule creates new file statistic by file size range module
func NewDetailFileModule(ctx *Context, order int, enabledRanges []int) Module {
	// Do nothing if verbose not enabled
	if len(enabledRanges) == 0 {
		return newVoidModule()
	}
	details := newDetailsFile(newRanges(), enabledRanges)
	handler := newDetailFileHandler(details)
	rend := newDetailFileRenderer(details, ctx.pd, order)

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
	hlers []scan.Handler
	rnd   []renderer
}

func (m *module) handlers() []scan.Handler {
	return m.hlers
}

func (m *module) renderers() []renderer {
	return m.rnd
}

// newVoidModule creates module that do nothing
func newVoidModule() Module {
	return &module{
		[]scan.Handler{},
		[]renderer{},
	}
}

func newModule(r renderer, handlers ...scan.Handler) Module {
	m := module{
		[]scan.Handler{},
		[]renderer{r},
	}
	m.hlers = append(m.hlers, handlers...)
	return &m
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

func render(renderers []renderer, p out.Printer) {
	sort.Slice(renderers, func(i, j int) bool {
		return renderers[i].order() < renderers[j].order()
	})

	for _, r := range renderers {
		r.render(p)
	}
}
