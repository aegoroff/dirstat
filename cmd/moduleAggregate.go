package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"text/tabwriter"
)

type moduleAggregate struct {
	aggregator map[string]countSizeAggregate
	total      *totalInfo
}

func (m *moduleAggregate) postScan() {

}

func (m *moduleAggregate) handler() fileHandler {
	return func(f *fileEntry) {
		ext := filepath.Ext(f.Name)
		a := m.aggregator[ext]
		a.Size += uint64(f.Size)
		a.Count++
		m.aggregator[ext] = a
	}
}

func (m *moduleAggregate) output(tw *tabwriter.Writer, w io.Writer) {
	m.total.CountFileExts = len(m.aggregator)

	extBySize := createSliceFromMap(m.aggregator, func(aggregate countSizeAggregate) int64 {
		return int64(aggregate.Size)
	})

	extByCount := createSliceFromMap(m.aggregator, func(aggregate countSizeAggregate) int64 {
		return aggregate.Count
	})

	sort.Sort(sort.Reverse(extBySize))
	sort.Sort(sort.Reverse(extByCount))

	const format = "%v\t%v\t%v\t%v\t%v\n"

	_, _ = fmt.Fprintf(w, "\nTOP %d file extensions by size:\n\n", Top)
	_, _ = fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

	outputTopTenExtensions(tw, extBySize, m.total, func(data containers, item *container) (int64, uint64) {
		count := m.aggregator[item.name].Count
		sz := uint64(item.size)
		return count, sz
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d file extensions by count:\n\n", Top)
	_, _ = fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

	outputTopTenExtensions(tw, extByCount, m.total, func(data containers, item *container) (int64, uint64) {
		count := item.size
		sz := m.aggregator[item.name].Size
		return count, sz
	})

	_ = tw.Flush()
}

func outputTopTenExtensions(tw *tabwriter.Writer, data containers, total *totalInfo, selector func(data containers, item *container) (int64, uint64)) {
	for i := 0; i < Top && i < len(data); i++ {
		h := data[i].name

		count, sz := selector(data, data[i])

		outputTopStatLine(tw, count, total, sz, h)
	}
}

func createSliceFromMap(sizeByExt map[string]countSizeAggregate, mapper func(countSizeAggregate) int64) containers {
	var result = make(containers, len(sizeByExt))
	i := 0
	for k, v := range sizeByExt {
		result[i] = &container{size: mapper(v), name: k}
		i++
	}
	return result
}
