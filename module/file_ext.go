package module

import (
	"dirstat/module/internal/sys"
	"path/filepath"
	"sort"
)

type extWorker struct {
	voidInit
	*fileFilter
	total      *totalInfo
	aggregator map[string]countSizeAggregate
}

type extRenderer struct {
	work *extWorker
	top  int
}

func newExtWorker(ctx *Context) *extWorker {
	w := extWorker{
		total:      ctx.total,
		aggregator: make(map[string]countSizeAggregate, 8192),
	}

	w.fileFilter = newFileFilter(w.onFile)

	return &w
}

func newExtRenderer(work *extWorker, top int) renderer {
	return &extRenderer{work: work, top: top}
}

// Worker methods

func (m *extWorker) finalize() {
	m.total.CountFileExts = len(m.aggregator)
}

func (m *extWorker) onFile(f *sys.FileEntry) {
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

	sizePrint := fileExtPrint{
		data:    extBySize,
		count:   func(f *file) int64 { return e.work.aggregator[f.path].Count },
		size:    func(f *file) uint64 { return uint64(f.size) },
		p:       p,
		headfmt: "\n<gray>TOP %d file extensions by size:</>\n\n",
		total:   e.work.total,
	}

	sizePrint.print(e.top)

	countPrint := fileExtPrint{
		data:    extByCount,
		count:   func(f *file) int64 { return f.size },
		size:    func(f *file) uint64 { return e.work.aggregator[f.path].Size },
		p:       p,
		headfmt: "\n<gray>TOP %d file extensions by count:</>\n\n",
		total:   e.work.total,
	}

	countPrint.print(e.top)
}

type fileExtPrint struct {
	data    files
	count   func(f *file) int64
	size    func(f *file) uint64
	p       printer
	headfmt string
	total   *totalInfo
}

func (fp *fileExtPrint) print(top int) {
	const format = "%v\t%v\t%v\t%v\t%v\t%v\n"

	fp.p.cprint(fp.headfmt, top)

	fp.printTableHead(format)

	fp.printTopTen(top, fp.data, func(data files, item *file) (int64, uint64) {
		count := fp.count(item)
		sz := fp.size(item)
		return count, sz
	})

	fp.p.flush()
}

func (fp *fileExtPrint) printTableHead(format string) {
	fp.p.tprint(format, " #", "Extension", "Count", "%", "Size", "%")
	fp.p.tprint(format, "--", "---------", "-----", "------", "----", "------")
}

func (fp *fileExtPrint) printTopTen(top int, data files, selector func(data files, item *file) (int64, uint64)) {
	for i := 0; i < top && i < len(data); i++ {
		h := data[i].path

		count, sz := selector(data, data[i])

		fp.total.printCountAndSizeStatLine(fp.p, i+1, count, sz, h)
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
