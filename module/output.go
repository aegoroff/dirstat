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

func outputTopStatLine(tw *tabwriter.Writer, count int64, total *totalInfo, sz uint64, title string) {
	percentOfCount := countPercent(count, total)
	percentOfSize := sizePercent(sz, total)

	_, _ = fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
}

func countPercent(count int64, total *totalInfo) float64 {
	return (float64(count) / float64(total.FilesTotal.Count)) * 100
}

func sizePercent(size uint64, total *totalInfo) float64 {
	return (float64(size) / float64(total.FilesTotal.Size)) * 100
}
