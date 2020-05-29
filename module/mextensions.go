package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"text/tabwriter"
)

// NewExtensionModule creates new file extensions statistic module
func NewExtensionModule(ctx *Context) Module {
	work := newExtWorker(ctx)
	rend := extRenderer{work}
	m := moduleExtensions{
		work,
		rend,
	}
	return &m
}

// NewExtensionHiddenModule creates new file extensions statistic module
// that has disabled output
func NewExtensionHiddenModule(ctx *Context) Module {
	work := newExtWorker(ctx)
	m := moduleExtensionsNoOut{
		work,
		emptyRenderer{},
	}
	return &m
}

type extWorker struct {
	total      *totalInfo
	aggregator map[string]countSizeAggregate
}

type extRenderer struct {
	extWorker
}

type moduleExtensions struct {
	extWorker
	extRenderer
}

type moduleExtensionsNoOut struct {
	extWorker
	emptyRenderer
}

func newExtWorker(ctx *Context) extWorker {
	return extWorker{
		total:      ctx.total,
		aggregator: make(map[string]countSizeAggregate),
	}
}

func (m *extWorker) init() {
}

func (m *extWorker) postScan() {
	m.total.CountFileExts = len(m.aggregator)
}

func (m *extWorker) folderHandler(_ *sys.FolderEntry) {

}

func (m *extWorker) fileHandler(f *sys.FileEntry) {
	ext := filepath.Ext(f.Path)
	a := m.aggregator[ext]
	a.Size += uint64(f.Size)
	a.Count++
	m.aggregator[ext] = a
}

func (m *moduleExtensions) output(tw *tabwriter.Writer, w io.Writer) {
	extBySize := createSliceFromMap(m.aggregator, func(aggregate countSizeAggregate) int64 {
		return int64(aggregate.Size)
	})

	extByCount := createSliceFromMap(m.aggregator, func(aggregate countSizeAggregate) int64 {
		return aggregate.Count
	})

	sort.Sort(sort.Reverse(extBySize))
	sort.Sort(sort.Reverse(extByCount))

	const format = "%v\t%v\t%v\t%v\t%v\n"

	_, _ = fmt.Fprintf(w, "\nTOP %d file extensions by size:\n\n", top)

	outputExtenstionsHead(tw, format)

	outputTopTenExtensions(tw, extBySize, m.total, func(data containers, item *container) (int64, uint64) {
		count := m.aggregator[item.name].Count
		sz := uint64(item.size)
		return count, sz
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d file extensions by count:\n\n", top)

	outputExtenstionsHead(tw, format)

	outputTopTenExtensions(tw, extByCount, m.total, func(data containers, item *container) (int64, uint64) {
		count := item.size
		sz := m.aggregator[item.name].Size
		return count, sz
	})

	_ = tw.Flush()
}

func outputExtenstionsHead(tw *tabwriter.Writer, format string) {
	_, _ = fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")
}

func outputTopTenExtensions(tw *tabwriter.Writer, data containers, total *totalInfo, selector func(data containers, item *container) (int64, uint64)) {
	for i := 0; i < top && i < len(data); i++ {
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
