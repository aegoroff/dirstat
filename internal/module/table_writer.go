package module

import (
	"fmt"
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
)

func newTableWriter(w out.Writable) table.Writer {
	tab := table.NewWriter()
	tab.SetAllowedRowLength(0)
	tab.SetOutputMirror(w.Writer())
	tab.SetStyle(table.StyleLight)
	tab.Style().Options.SeparateColumns = true
	tab.Style().Options.DrawBorder = true
	return tab
}

func appendHeaders(heads []string, tab table.Writer) {
	headers := table.Row{}
	for _, head := range heads {
		headers = append(headers, head)
	}
	tab.AppendHeader(headers)
}

func percentTransformer(val interface{}) string {
	v := val.(float64)
	if v >= 90.0 {
		return color.Sprintf("<red>%.2f%%</>", v)
	}
	if v >= 70.0 {
		return color.Sprintf("<yellow>%.2f%%</>", v)
	}
	return fmt.Sprintf("%.2f%%", v)
}

func sizeTransformer(val interface{}) string {
	sz := val.(uint64)
	return humanize.IBytes(sz)
}
