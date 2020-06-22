package module

import (
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"text/template"
)

type totalRenderer struct {
	total *totalInfo
}

func newTotalRenderer(ctx *Context) renderer {
	return &totalRenderer{ctx.total}
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
