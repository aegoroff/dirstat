package module

import (
	"dirstat/module/internal/sys"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"path/filepath"
	"text/template"
)

type totalFileHandler struct {
	total *totalInfo
}

type totalFolderHandler struct {
	total *totalInfo
}

type totalRenderer struct {
	*baseRenderer
	total *totalInfo
}

func newTotalFileHandler(ctx *Context) sys.Handler {
	w := totalFileHandler{
		total: ctx.total,
	}

	return &w
}

func newTotalFolderHandler(ctx *Context) sys.Handler {
	w := totalFolderHandler{
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

func (m *totalFileHandler) Handle(evt *sys.ScanEvent) {
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

func (m *totalFolderHandler) Handle(*sys.ScanEvent) {
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
