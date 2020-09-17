package module

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"io"
	"text/tabwriter"
)

type printer interface {
	writer() io.Writer
	twriter() *tabwriter.Writer
	flush()

	// tprint prints using tab writer
	tprint(format string, a ...interface{})

	// cprint prints data with suppport colorizing
	cprint(format string, a ...interface{})
}

type prn struct {
	tw *tabwriter.Writer
	w  io.Writer
}

func (r *prn) writer() io.Writer {
	return r.w
}

func (r *prn) twriter() *tabwriter.Writer {
	return r.tw
}

func (r *prn) flush() {
	_ = r.tw.Flush()
}

func (r *prn) tprint(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(r.tw, format, a...)
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
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	p := prn{
		tw: tw,
		w:  w,
	}
	return &p
}
