package module

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"text/tabwriter"
)

type renderContext interface {
	writer() io.Writer
	flush()
	write(format string, a ...interface{})
	writetab(format string, a ...interface{})
}

type rctx struct {
	tabw *tabwriter.Writer
	wr   io.Writer
}

func (r *rctx) writer() io.Writer {
	return r.wr
}

func (r *rctx) flush() {
	_ = r.tabw.Flush()
}

func (r *rctx) write(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(r.wr, format, a...)
}

func (r *rctx) writetab(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(r.tabw, format, a...)
}

func render(w io.Writer, renderers []renderer) {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	c := rctx{
		tabw: tw,
		wr:   w,
	}

	for _, r := range renderers {
		r.output(&c)
	}
}

func (t *totalInfo) outputTopStatLine(rc renderContext, count int64, sz uint64, title string) {
	percentOfCount := t.countPercent(count)
	percentOfSize := t.sizePercent(sz)

	rc.writetab("%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
}

func (t *totalInfo) countPercent(count int64) float64 {
	return (float64(count) / float64(t.FilesTotal.Count)) * 100
}

func (t *totalInfo) sizePercent(size uint64) float64 {
	return (float64(size) / float64(t.FilesTotal.Size)) * 100
}
