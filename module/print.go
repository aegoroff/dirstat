package module

import (
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"io"
)

type printer interface {
	writer() io.Writer

	// cprint prints data with suppport colorizing
	cprint(format string, a ...interface{})
}

type prn struct {
	w io.Writer
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
