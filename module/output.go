package module

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"text/tabwriter"
)

func render(w io.Writer, renderers []renderer) {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	for _, r := range renderers {
		r.output(tw, w)
	}
}

func (t *totalInfo) outputTopStatLine(tw *tabwriter.Writer, count int64, sz uint64, title string) {
	percentOfCount := t.countPercent(count)
	percentOfSize := t.sizePercent(sz)

	_, _ = fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
}

func (t *totalInfo) countPercent(count int64) float64 {
	return (float64(count) / float64(t.FilesTotal.Count)) * 100
}

func (t *totalInfo) sizePercent(size uint64) float64 {
	return (float64(size) / float64(t.FilesTotal.Size)) * 100
}
