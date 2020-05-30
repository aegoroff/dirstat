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

func (m *extWorker) folderHandler(*sys.FolderEntry) {

}

func (m *extWorker) fileHandler(f *sys.FileEntry) {
	ext := filepath.Ext(f.Path)
	a := m.aggregator[ext]
	a.Size += uint64(f.Size)
	a.Count++
	m.aggregator[ext] = a
}

func (e *extRenderer) output(p printer) {
	extBySize := e.evolventMap(func(agr countSizeAggregate) int64 {
		return int64(agr.Size)
	})

	extByCount := e.evolventMap(func(agr countSizeAggregate) int64 {
		return agr.Count
	})

	sort.Sort(sort.Reverse(extBySize))
	sort.Sort(sort.Reverse(extByCount))

	const format = "%v\t%v\t%v\t%v\t%v\n"

	p.print("\nTOP %d file extensions by size:\n\n", top)

	e.outputTableHead(p, format)

	e.outputTopTen(p, extBySize, func(data containers, item *container) (int64, uint64) {
		count := e.aggregator[item.name].Count
		sz := uint64(item.size)
		return count, sz
	})

	p.flush()

	p.print("\nTOP %d file extensions by count:\n\n", top)

	e.outputTableHead(p, format)

	e.outputTopTen(p, extByCount, func(data containers, item *container) (int64, uint64) {
		count := item.size
		sz := e.aggregator[item.name].Size
		return count, sz
	})

	p.flush()
}

func (e *extRenderer) outputTableHead(p printer, format string) {
	p.printtab(format, "Extension", "Count", "%", "Size", "%")
	p.printtab(format, "---------", "-----", "------", "----", "------")
}

func (e *extRenderer) outputTopTen(p printer, data containers, selector func(data containers, item *container) (int64, uint64)) {
	for i := 0; i < top && i < len(data); i++ {
		h := data[i].name

		count, sz := selector(data, data[i])

		e.total.printTopStatLine(p, count, sz, h)
	}
}

func (e *extRenderer) evolventMap(mapper func(countSizeAggregate) int64) containers {
	var result = make(containers, len(e.aggregator))
	i := 0
	for k, v := range e.aggregator {
		result[i] = &container{size: mapper(v), name: k}
		i++
	}
	return result
}
