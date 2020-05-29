package module

import (
	"dirstat/module/internal/sys"
	"path/filepath"
	"sort"
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
func NewExtensionHiddenModule(rc *Context) Module {
	work := newExtWorker(rc)
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

func newExtWorker(rc *Context) extWorker {
	return extWorker{
		total:      rc.total,
		aggregator: make(map[string]countSizeAggregate),
	}
}

func (m *extWorker) init() {
}

func (m *extWorker) postScan() {
	m.total.CountFileExts = len(m.aggregator)
}

func (m *extWorker) folderHandler(*sys.FolderEntry) {

}

func (m *extWorker) fileHandler(f *sys.FileEntry) {
	ext := filepath.Ext(f.Path)
	a := m.aggregator[ext]
	a.Size += uint64(f.Size)
	a.Count++
	m.aggregator[ext] = a
}

func (e *extRenderer) output(rc renderContext) {
	extBySize := createSliceFromMap(e.aggregator, func(aggregate countSizeAggregate) int64 {
		return int64(aggregate.Size)
	})

	extByCount := createSliceFromMap(e.aggregator, func(aggregate countSizeAggregate) int64 {
		return aggregate.Count
	})

	sort.Sort(sort.Reverse(extBySize))
	sort.Sort(sort.Reverse(extByCount))

	const format = "%v\t%v\t%v\t%v\t%v\n"

	rc.write("\nTOP %d file extensions by size:\n\n", top)

	e.outputTableHead(rc, format)

	e.outputTopTen(rc, extBySize, func(data containers, item *container) (int64, uint64) {
		count := e.aggregator[item.name].Count
		sz := uint64(item.size)
		return count, sz
	})

	rc.flush()

	rc.write("\nTOP %d file extensions by count:\n\n", top)

	e.outputTableHead(rc, format)

	e.outputTopTen(rc, extByCount, func(data containers, item *container) (int64, uint64) {
		count := item.size
		sz := e.aggregator[item.name].Size
		return count, sz
	})

	rc.flush()
}

func (e *extRenderer) outputTableHead(rc renderContext, format string) {
	rc.writetab(format, "Extension", "Count", "%", "Size", "%")
	rc.writetab(format, "---------", "-----", "------", "----", "------")
}

func (e *extRenderer) outputTopTen(rc renderContext, data containers, selector func(data containers, item *container) (int64, uint64)) {
	for i := 0; i < top && i < len(data); i++ {
		h := data[i].name

		count, sz := selector(data, data[i])

		e.total.outputTopStatLine(rc, count, sz, h)
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
