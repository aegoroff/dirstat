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
	fileHandler(f *sys.FileEntry)
	folderHandler(f *sys.FolderEntry)
	output(tw *tabwriter.Writer, w io.Writer)
	postScan()
	init()
}

// Context defines modules context
type Context struct {
	total          *totalInfo
	rangeAggregate map[Range]fileStat
}

// Execute runs modules over path specified
func Execute(path string, fs afero.Fs, w io.Writer, modules []Module) {
	var handlers []sys.ScanHandler
	for _, m := range modules {
		m.init()
		handlers = append(handlers, scanEventHandler(m))
	}
	sys.Scan(path, fs, handlers)

	for _, m := range modules {
		m.postScan()
	}

	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	for _, m := range modules {
		m.output(tw, w)
	}
}

func scanEventHandler(m Module) sys.ScanHandler {
	return func(evt *sys.ScanEvent) {
		if evt.Folder != nil {
			m.folderHandler(evt.Folder)
		} else if evt.File != nil {
			m.fileHandler(evt.File)
		}
	}
}

// NewContext create new module's context that needed to create new modules
func NewContext() *Context {
	total := totalInfo{}

	ctx := Context{
		total:          &total,
		rangeAggregate: make(map[Range]fileStat),
	}
	return &ctx
}

// NewFoldersModule creates new folders module
func NewFoldersModule(ctx *Context) Module {
	m := moduleFolders{
		ctx.total,
		rbtree.NewRbTree(),
		rbtree.NewRbTree(),
	}
	return &m
}

// NewFoldersHiddenModule creates new folders module
// that has disabled output
func NewFoldersHiddenModule(ctx *Context) Module {
	m := moduleFolders{
		ctx.total,
		rbtree.NewRbTree(),
		rbtree.NewRbTree(),
	}
	h := moduleFoldersNoOut{
		m,
	}
	return &h
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
