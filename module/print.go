package module

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"text/tabwriter"
)

type printer interface {
	writer() io.Writer
	flush()
	print(format string, a ...interface{})
}

type prn struct {
	tw *tabwriter.Writer
	w  io.Writer
}

func (r *prn) writer() io.Writer {
	return r.w
}

func (r *prn) flush() {
	_ = r.tw.Flush()
}

func (r *prn) print(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(r.tw, format, a...)
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
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	p := prn{
		tw: tw,
		w:  w,
	}
	return &p
}
