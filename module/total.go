package module

import (
	"dirstat/module/internal/sys"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"path/filepath"
	"text/template"
)

type totalFileWorker struct {
	total *totalInfo
}

type totalFolderWorker struct {
	total *totalInfo
}

type totalRenderer struct {
	*baseRenderer
	total *totalInfo
}

func newTotalFileWorker(ctx *Context) worker {
	w := totalFileWorker{
		total: ctx.total,
	}

	return &w
}

func newTotalFolderWorker(ctx *Context) worker {
	w := totalFolderWorker{
		total: ctx.total,
	}

	return &w
}

func newTotalRenderer(ctx *Context, order int) renderer {
	return &totalRenderer{
		baseRenderer: newBaseRenderer(order),
		total:        ctx.total,
	}
}

// Worker methods

func (m *totalFileWorker) handler(evt *sys.ScanEvent) {
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

func (m *totalFolderWorker) handler(*sys.ScanEvent) {
	m.total.CountFolders++
}

// Renderer method

func (m *totalRenderer) print(p printer) {
	m.total.countExtensions()
	const totalTemplate = `
Total files:            {{.FilesTotal.Count}} ({{.FilesTotal.Size | toBytesString }})
Total folders:          {{.CountFolders}}
Total file extensions:  {{.CountFileExts}}`

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))

	_, _ = color.Set(color.FgGray)
	_ = report.Execute(p.writer(), m.total)

	_, _ = color.Reset()
}
