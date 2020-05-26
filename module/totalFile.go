package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"io"
	"text/tabwriter"
)

type moduleTotalFile struct {
	moduleTotal
	aggregate map[Range]fileStat
}

func (m *moduleTotalFile) postScan() {

}

func (m *moduleTotalFile) handler() sys.FileHandler {
	return func(f *sys.FileEntry) {
	}
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
