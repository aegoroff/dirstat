package module

import (
	"dirstat/module/internal/sys"
	"github.com/dustin/go-humanize"
	"text/template"
	"time"
)

// NewTotalModule creates new total statistic module
func NewTotalModule(ctx *Context) Module {
	m := moduleTotal{
		start: time.Now(),
		total: ctx.total,
	}
	return &m
}

type moduleTotal struct {
	total *totalInfo
	start time.Time
}

func (m *moduleTotal) init() {
}

func (m *moduleTotal) postScan() {
	m.total.ReadingTime = time.Since(m.start)
}

func (m *moduleTotal) folderHandler(*sys.FolderEntry) {

}

func (m *moduleTotal) fileHandler(f *sys.FileEntry) {
	// Accumulate file statistic
	m.total.FilesTotal.Count++
	m.total.FilesTotal.Size += uint64(f.Size)
}

func (m *moduleTotal) print(p printer) {
	const totalTemplate = `
Total files:            {{.FilesTotal.Count}} ({{.FilesTotal.Size | toBytesString }})
Total folders:          {{.CountFolders}}
Total file extensions:  {{.CountFileExts}}

Read taken:    {{.ReadingTime}}
`

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
	_ = report.Execute(p.writer(), m.total)
}
