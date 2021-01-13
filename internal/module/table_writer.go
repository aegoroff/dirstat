package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
)

type tableWriter struct {
	p   out.Printer
	tab table.Writer
}

func newTableWriter(p out.Printer) *tableWriter {
	tab := table.NewWriter()
	tab.SetAllowedRowLength(0)
	tab.SetOutputMirror(p.Writer())
	tab.SetStyle(table.StyleLight)
	tab.Style().Options.SeparateColumns = true
	tab.Style().Options.DrawBorder = true

	return &tableWriter{p: p, tab: tab}
}

func (t *tableWriter) appendHeaders(heads []string) {
	headers := table.Row{}
	for _, head := range heads {
		headers = append(headers, head)
	}
	t.tab.AppendHeader(headers)
}

func (t *tableWriter) percentTransformer(val interface{}) string {
	v := val.(float64)
	if v >= 90.0 {
		return t.p.Sprintf("<red>%.2f%%</>", v)
	}
	if v >= 70.0 {
		return t.p.Sprintf("<yellow>%.2f%%</>", v)
	}
	return t.p.Sprintf("%.2f%%", v)
}

func (*tableWriter) sizeTransformer(val interface{}) string {
	sz := val.(uint64)
	return humanize.IBytes(sz)
}
