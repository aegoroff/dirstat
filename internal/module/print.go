package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
	"sort"
)

type printer interface {
	out.Printer
	createTab() table.Writer
}

type prn struct {
	p out.Printer
}

func (r *prn) Cprint(format string, a ...interface{}) {
	r.p.Cprint(format, a...)
}

func (r *prn) Println() {
	r.p.Println()
}

func (r *prn) Writer() io.WriteCloser {
	return r.p.Writer()
}

func (r *prn) createTab() table.Writer {
	tab := table.NewWriter()
	tab.SetAllowedRowLength(0)
	tab.SetOutputMirror(r.Writer())
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

func human(n int64) string {
	return humanize.IBytes(uint64(n))
}

func render(p out.Printer, renderers []renderer) {
	pn := newPrinter(p)

	sort.Slice(renderers, func(i, j int) bool {
		return renderers[i].order() < renderers[j].order()
	})

	for _, r := range renderers {
		r.print(pn)
	}
}

func newPrinter(p out.Printer) printer {
	return &prn{
		p: p,
	}
}
