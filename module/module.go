package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"io"
	"text/tabwriter"
)

// Module defines working modules interface
type Module interface {
	worker
	renderer
}

type worker interface {
	fileHandler(f *sys.FileEntry)
	folderHandler(f *sys.FolderEntry)
	postScan()
	init()
}

type renderer interface {
	output(tw *tabwriter.Writer, w io.Writer)
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
