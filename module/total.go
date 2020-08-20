package module

import (
	"dirstat/module/internal/sys"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"path/filepath"
	"text/template"
)

type totalWorker struct {
	voidInit
	total *totalInfo
}

type totalRenderer struct {
	total *totalInfo
}

func newTotalWorker(ctx *Context) *totalWorker {
	w := totalWorker{
		total: ctx.total,
	}

	return &w
}

func newTotalRenderer(ctx *Context) renderer {
	return &totalRenderer{ctx.total}
}

// Worker methods

func (m *totalWorker) finalize() {
	m.total.CountFileExts = len(m.total.extensions)
}

func (m *totalWorker) handler(evt *sys.ScanEvent) {
	if evt.Folder != nil {
		m.total.CountFolders++
	} else if evt.File != nil {
		f := evt.File
		// Accumulate file statistic
		m.total.FilesTotal.Count++
		m.total.FilesTotal.Size += uint64(f.Size)

		// Accumulate file extensions statistic
		ext := filepath.Ext(f.Path)
		a := m.total.extensions[ext]
		a.Size += uint64(f.Size)
		a.Count++
		m.total.extensions[ext] = a
	}
}

// Renderer method

func (m *totalRenderer) print(p printer) {
	const totalTemplate = `
Total files:            {{.FilesTotal.Count}} ({{.FilesTotal.Size | toBytesString }})
Total folders:          {{.CountFolders}}
Total file extensions:  {{.CountFileExts}}`

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))

	_, _ = color.Set(color.FgGray)
	_ = report.Execute(p.writer(), m.total)

	_, _ = color.Reset()
}
