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
	m.total.FilesTotal++
	m.total.FilesSize += f.Size

	// Accumulate file extensions statistic

	fn := &folder{
		path:  filepath.Ext(f.Path),
		count: 1,
		size:  f.Size,
	}
	n, ok := m.total.extensions.Search(fn)
	if ok {
		fold := n.(*folder)
		fold.size += fn.size
		fold.count++
	} else {
		m.total.extensions.Insert(fn)
	}
}

func (m *totalFolderHandler) Handle(*scan.Event) {
	m.total.FoldersTotal++
}

// Renderer method

func (m *totalRenderer) render(p out.Printer) {
	m.total.ExtensionsTotal = m.total.extensions.Len()
	const totalTemplate = `
Total files:            {{.FilesTotal | humanCount}} ({{.FilesSize | humanSize}})
Total folders:          {{.FoldersTotal | humanCount}}
Total file extensions:  {{.ExtensionsTotal | humanCount}}`

	transform := template.FuncMap{
		"humanSize":  humanSize,
		"humanCount": humanize.Comma,
	}
	tpl := template.New("totalstat").Funcs(transform)

	var report = template.Must(tpl.Parse(totalTemplate))

	_, _ = color.Set(color.FgGray)
	_ = report.Execute(p.Writer(), m.total)

	_, _ = color.Reset()
}
