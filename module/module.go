package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"io"
	"text/tabwriter"
	"time"
)

// Module defines working modules interface
type Module interface {
	handler() sys.FileHandler
	output(tw *tabwriter.Writer, w io.Writer)
	postScan()
	init()
}

// Context defines modules context
type Context struct {
	total          *totalInfo
	folders        map[string]*container
	rangeAggregate map[Range]fileStat
}

// Execute runs modules over path specified
func Execute(path string, fs afero.Fs, w io.Writer, ctx *Context, modules []Module) {
	foldersHandler := func(fe *sys.FolderEntry) {
		ctx.folders[fe.Name] = &container{name: fe.Name, count: fe.Count, size: fe.Size}
		ctx.total.CountFolders++
	}
	executeModules(path, fs, w, foldersHandler, modules)
}

// NewContext create new module's context that needed to create new modules
func NewContext() *Context {
	total := totalInfo{}
	folders := make(map[string]*container)

	ctx := Context{
		total:          &total,
		folders:        folders,
		rangeAggregate: make(map[Range]fileStat),
	}
	return &ctx
}

// NewFoldersModule creates new folders module
func NewFoldersModule(ctx *Context) Module {
	m := moduleFolders{
		ctx.total,
		ctx.folders,
		rbtree.NewRbTree(),
	}
	return &m
}

// NewTotalModule creates new total statistic module
func NewTotalModule(ctx *Context) Module {
	m := moduleTotal{
		start: time.Now(),
		total: ctx.total,
	}
	return &m
}

// NewTotalFileModule creates new total file statistic module
func NewTotalFileModule(ctx *Context) Module {
	m := moduleTotalFile{
		total:     ctx.total,
		aggregate: ctx.rangeAggregate,
	}
	return &m
}

// NewRangeModule creates new file statistic by file size range module
func NewRangeModule(ctx *Context, verbose bool, enabledRanges []int) Module {
	m := moduleRange{
		verbose:       verbose,
		enabledRanges: enabledRanges,
		aggregate:     ctx.rangeAggregate,
		distribution:  make(map[Range]containers),
	}
	return &m
}

// NewRangeHiddenModule creates new file statistic by file size range module
// that has disabled output
func NewRangeHiddenModule(ctx *Context) Module {
	m := moduleRange{
		verbose:       false,
		enabledRanges: []int{},
		aggregate:     ctx.rangeAggregate,
		distribution:  make(map[Range]containers),
	}
	h := moduleRangeNoOut{
		m,
	}
	return &h
}

// NewExtensionModule creates new file extensions statistic module
func NewExtensionModule(ctx *Context) Module {
	m := moduleExtensions{
		total:      ctx.total,
		aggregator: make(map[string]countSizeAggregate),
	}
	return &m
}

// NewExtensionHiddenModule creates new file extensions statistic module
// that has disabled output
func NewExtensionHiddenModule(ctx *Context) Module {
	m := moduleExtensions{
		total:      ctx.total,
		aggregator: make(map[string]countSizeAggregate),
	}
	h := moduleExtensionsNoOut{
		m,
	}
	return &h
}

// NewTopFilesModule creates new top files statistic module
func NewTopFilesModule(_ *Context) Module {
	m := moduleTopFiles{
		tree: rbtree.NewRbTree(),
	}
	return &m
}

// NewTopFilesHiddenModule creates new top files statistic module
// that has disabled output
func NewTopFilesHiddenModule(_ *Context) Module {
	m := moduleTopFiles{
		tree: rbtree.NewRbTree(),
	}
	h := moduleTopFilesNoOut{
		m,
	}
	return &h
}

func executeModules(path string, fs afero.Fs, w io.Writer, fh sys.FolderHandler, modules []Module) {
	var handlers []sys.FileHandler
	for _, m := range modules {
		m.init()
		handlers = append(handlers, m.handler())
	}
	sys.Scan(path, fs, fh, handlers)

	for _, m := range modules {
		m.postScan()
	}

	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	for _, m := range modules {
		m.output(tw, w)
	}
}

func outputTopStatLine(tw *tabwriter.Writer, count int64, total *totalInfo, sz uint64, title string) {
	percentOfCount := countPercent(count, total)
	percentOfSize := sizePercent(sz, total)

	_, _ = fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
}

func countPercent(count int64, total *totalInfo) float64 {
	return (float64(count) / float64(total.FilesTotal.Count)) * 100
}

func sizePercent(size uint64, total *totalInfo) float64 {
	return (float64(size) / float64(total.FilesTotal.Size)) * 100
}
