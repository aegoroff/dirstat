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

func (t *tableWriter) addHeaders(heads []string) {
	headers := table.Row{}
	for _, head := range heads {
		headers = append(headers, head)
	}
	t.tab.AppendHeader(headers)
}

func (t *tableWriter) configColumns(cc []table.ColumnConfig) {
	t.tab.SetColumnConfigs(cc)
}

func (t *tableWriter) render() string {
	return t.tab.Render()
}

func (t *tableWriter) addRow(row table.Row, configs ...table.RowConfig) {
	t.tab.AppendRow(row, configs...)
}

func (t *tableWriter) percentTransformer(val any) string {
	v := val.(float64)
	if v >= 90.0 {
		return t.p.Sprintf("<red>%.2f%%</>", v)
	}
	if v >= 70.0 {
		return t.p.Sprintf("<yellow>%.2f%%</>", v)
	}
	return t.p.Sprintf("%.2f%%", v)
}

func (*tableWriter) sizeTransformer(val any) string {
	szu, ok := val.(uint64)
	if ok {
		return humanize.IBytes(szu)
	}
	szi, ok := val.(int64)
	if ok {
		return humanize.IBytes(uint64(szi))
	}
	return ""
}

func (*tableWriter) countTransformer(val any) string {
	sz := val.(int64)
	return humanize.Comma(sz)
}
