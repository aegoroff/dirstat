package module

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// NewTotalFileModule creates new total file statistic module
func NewTotalFileModule(ctx *Context) Module {
	m := moduleTotalFile{
		total:     ctx.total,
		aggregate: ctx.rangeAggregate,
	}
	return &m
}

type moduleTotalFile struct {
	emptyWorker
	total     *totalInfo
	aggregate map[Range]fileStat
}

func (m *moduleTotalFile) output(tw *tabwriter.Writer, w io.Writer) {
	_, _ = fmt.Fprintf(w, "Total files stat:\n\n")

	const format = "%v\t%v\t%v\t%v\t%v\n"

	_, _ = fmt.Fprintf(tw, format, "File size", "Amount", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "------", "------", "----", "------")

	heads := createRangesHeads()
	for i, r := range fileSizeRanges {
		count := m.aggregate[r].TotalFilesCount
		sz := m.aggregate[r].TotalFilesSize

		outputTopStatLine(tw, count, m.total, sz, heads[i])
	}
	_ = tw.Flush()
}
