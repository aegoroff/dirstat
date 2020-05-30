package module

import (
	"dirstat/module/internal/sys"
	"github.com/dustin/go-humanize"
	"text/template"
	"time"
)

// NewTotalModule creates new total statistic module
func NewTotalModule(ctx *Context) Module {
	work := newTotalWorker(ctx)
	rend := newTotalRenderer(work)

	m := NewModule()
	m.addWorker(work)
	m.addRenderer(rend)
	return m
}

func newTotalRenderer(work *totalWorker) renderer {
	return &totalRenderer{work}
}

func newTotalWorker(ctx *Context) *totalWorker {
	return &totalWorker{
		start: time.Now(),
		total: ctx.total,
	}
}

type totalWorker struct {
	total *totalInfo
	start time.Time
}

type totalRenderer struct {
	work *totalWorker
}

func (m *totalWorker) init() {
}

func (m *totalWorker) finalize() {
	m.total.ReadingTime = time.Since(m.start)
}

func (m *totalWorker) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}

	// Accumulate file statistic
	m.total.FilesTotal.Count++
	m.total.FilesTotal.Size += uint64(evt.File.Size)
}

func (m *totalRenderer) print(p printer) {
	const totalTemplate = `
Total files:            {{.FilesTotal.Count}} ({{.FilesTotal.Size | toBytesString }})
Total folders:          {{.CountFolders}}
Total file extensions:  {{.CountFileExts}}

Read taken:    {{.ReadingTime}}
`

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
	_ = report.Execute(p.writer(), m.work.total)
}
