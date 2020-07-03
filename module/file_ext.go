package module

import (
	"dirstat/module/internal/sys"
	"path/filepath"
	"sort"
)

type extWorker struct {
	total      *totalInfo
	aggregator map[string]countSizeAggregate
}

type extRenderer struct {
	work *extWorker
	top  int
}

func newExtWorker(ctx *Context) *extWorker {
	return &extWorker{
		total:      ctx.total,
		aggregator: make(map[string]countSizeAggregate),
	}
}

func newExtRenderer(work *extWorker, top int) renderer {
	return &extRenderer{work: work, top: top}
}

// Worker methods

func (*extWorker) init() {
}

func (m *extWorker) finalize() {
	m.total.CountFileExts = len(m.aggregator)
}

func (m *extWorker) handler(evt *sys.ScanEvent) {
	f := evt.File

	// Accumulate file statistic
	m.total.FilesTotal.Count++
	m.total.FilesTotal.Size += uint64(f.Size)

	ext := filepath.Ext(f.Path)
	a := m.aggregator[ext]
	a.Size += uint64(f.Size)
	a.Count++
	m.aggregator[ext] = a
}

// Renderer method

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

	p.cprint("\n<gray>TOP %d file extensions by size:</>\n\n", e.top)

	e.printTableHead(p, format)

	e.printTopTen(p, extBySize, func(data files, item *file) (int64, uint64) {
		count := e.work.aggregator[item.path].Count
		sz := uint64(item.size)
		return count, sz
	})

	p.flush()

	p.cprint("\n<gray>TOP %d file extensions by count:</>\n\n", e.top)

	e.printTableHead(p, format)

	e.printTopTen(p, extByCount, func(data files, item *file) (int64, uint64) {
		count := item.size
		sz := e.work.aggregator[item.path].Size
		return count, sz
	})

	p.flush()
}

func (*extRenderer) printTableHead(p printer, format string) {
	p.print(format, "Extension", "Count", "%", "Size", "%")
	p.print(format, "---------", "-----", "------", "----", "------")
}

func (e *extRenderer) printTopTen(p printer, data files, selector func(data files, item *file) (int64, uint64)) {
	for i := 0; i < e.top && i < len(data); i++ {
		h := data[i].path

		count, sz := selector(data, data[i])

		e.work.total.printCountAndSizeStatLine(p, count, sz, h)
	}
}

func (e *extRenderer) evolventMap(mapper func(countSizeAggregate) int64) files {
	var result = make(files, len(e.work.aggregator))
	i := 0
	for k, v := range e.work.aggregator {
		result[i] = &file{size: mapper(v), path: k}
		i++
	}
	return result
}
