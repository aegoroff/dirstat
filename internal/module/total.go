package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
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

func newTotalFileHandler(ctx *Context) scan.Handler {
	w := totalFileHandler{
		total: ctx.total,
	}

	return newOnlyFilesHandler(&w)
}

func newTotalFolderHandler(ctx *Context) scan.Handler {
	w := totalFolderHandler{
		total: ctx.total,
	}

	return newOnlyFoldersHandler(&w)
}

func newTotalRenderer(ctx *Context, order int) renderer {
	return &totalRenderer{
		baseRenderer: newBaseRenderer(order),
		total:        ctx.total,
	}
}

// Handler methods

func (m *totalFileHandler) Handle(evt *scan.Event) {
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

func (m *totalFolderHandler) Handle(*scan.Event) {
	m.total.CountFolders++
}

// Renderer method

func (m *totalRenderer) render(p out.Printer) {
	m.total.countExtensions()
	const totalTemplate = `
Total files:            {{.FilesTotal.Count | humanCount}} ({{.FilesTotal.Size | humanSize}})
Total folders:          {{.CountFolders | humanCount}}
Total file extensions:  {{.CountFileExts | humanCountInt}}`

	transform := template.FuncMap{
		"humanSize":     humanize.IBytes,
		"humanCount":    humanize.Comma,
		"humanCountInt": func(i int) string { return humanize.Comma(int64(i)) },
	}
	tpl := template.New("totalstat").Funcs(transform)

	var report = template.Must(tpl.Parse(totalTemplate))

	_, _ = color.Set(color.FgGray)
	_ = report.Execute(p.Writer(), m.total)

	_, _ = color.Reset()
}
