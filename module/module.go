package module

import (
	"dirstat/module/internal/sys"
	"github.com/spf13/afero"
	"io"
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
	output(p printer)
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
	var handlers []sys.ScanHandler
	var renderers []renderer

	for _, m := range modules {
		renderers = append(renderers, m)
		m.init()
		handlers = append(handlers, scanEventHandler(m))
	}

	sys.Scan(path, fs, handlers)

	for _, m := range modules {
		m.postScan()
	}

	render(w, renderers)
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
