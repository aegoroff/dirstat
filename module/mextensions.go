package module

import (
	"dirstat/module/internal/sys"
	"path/filepath"
	"sort"
)

// NewExtensionModule creates new file extensions statistic module
func NewExtensionModule(ctx *Context) Module {
	work := newExtWorker(ctx)
	rend := &extRenderer{*work}
	m := module{
		[]worker{work},
		[]renderer{rend},
	}
	return &m
}

// NewExtensionHiddenModule creates new file extensions statistic module
// that has disabled output
func NewExtensionHiddenModule(ctx *Context) Module {
	work := newExtWorker(ctx)
	m := module{
		[]worker{work},
		[]renderer{},
	}
	return &m
}

type extWorker struct {
	total      *totalInfo
	aggregator map[string]countSizeAggregate
	top        int
}

type extRenderer struct {
	extWorker
}

func newExtWorker(ctx *Context) *extWorker {
	return &extWorker{
		total:      ctx.total,
		aggregator: make(map[string]countSizeAggregate),
		top:        ctx.top,
	}
}

func (m *extWorker) init() {
}

func (m *extWorker) finalize() {
	m.total.CountFileExts = len(m.aggregator)
}

func (m *extWorker) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}
	f := evt.File

	ext := filepath.Ext(f.Path)
	a := m.aggregator[ext]
	a.Size += uint64(f.Size)
	a.Count++
	m.aggregator[ext] = a
}

func (e *extRenderer) print(p printer) {
	extBySize := e.evolventMap(func(agr countSizeAggregate) int64 {
		return int64(agr.Size)
	})

	extByCount := e.evolventMap(func(agr countSizeAggregate) int64 {
		return agr.Count
	})

	sort.Sort(sort.Reverse(extBySize))
	sort.Sort(sort.Reverse(extByCount))

	const format = "%v\t%v\t%v\t%v\t%v\n"

	p.print("\nTOP %d file extensions by size:\n\n", e.top)

	e.printTableHead(p, format)

	e.printTopTen(p, extBySize, func(data containers, item *container) (int64, uint64) {
		count := e.aggregator[item.name].Count
		sz := uint64(item.size)
		return count, sz
	})

	p.flush()

	p.print("\nTOP %d file extensions by count:\n\n", e.top)

	e.printTableHead(p, format)

	e.printTopTen(p, extByCount, func(data containers, item *container) (int64, uint64) {
		count := item.size
		sz := e.aggregator[item.name].Size
		return count, sz
	})

	p.flush()
}

func (e *extRenderer) printTableHead(p printer, format string) {
	p.printtab(format, "Extension", "Count", "%", "Size", "%")
	p.printtab(format, "---------", "-----", "------", "----", "------")
}

func (e *extRenderer) printTopTen(p printer, data containers, selector func(data containers, item *container) (int64, uint64)) {
	for i := 0; i < e.top && i < len(data); i++ {
		h := data[i].name

		count, sz := selector(data, data[i])

		e.total.printCountAndSizeStatLine(p, count, sz, h)
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
