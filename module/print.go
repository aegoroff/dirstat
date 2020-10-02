package module

import (
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
)

type printer interface {
	writer() io.Writer
	createTab() table.Writer

	// cprint prints data with suppport colorizing
	cprint(format string, a ...interface{})
}

type prn struct {
	w io.Writer
}

func (r *prn) createTab() table.Writer {
	tab := table.NewWriter()
	tab.SetAllowedRowLength(0)
	tab.SetOutputMirror(r.w)
	tab.SetStyle(table.StyleLight)
	tab.Style().Options.SeparateColumns = true
	tab.Style().Options.DrawBorder = true
	return tab
}

func (r *prn) writer() io.Writer {
	return r.w
}

func (r *prn) cprint(format string, a ...interface{}) {
	color.Fprintf(r.w, format, a...)
}

func human(n int64) string {
	return humanize.IBytes(uint64(n))
}

func render(w io.Writer, renderers []renderer) {
	p := newPrinter(w)

	for _, r := range renderers {
		r.print(p)
	}
}

func newPrinter(w io.Writer) printer {
	p := prn{
		w: w,
	}
	return &p
}
