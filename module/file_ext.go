package module

import (
	"sort"
)

type extRenderer struct {
	total *totalInfo
	top   int
}

func newExtRenderer(ctx *Context) renderer {
	return &extRenderer{total: ctx.total, top: ctx.top}
}

// Renderer method

func (e *extRenderer) print(p printer) {
	extBySize := e.evolventMap(func(agr countSizeAggregate) int64 {
		return int64(agr.Size)
	})

	extByCount := e.evolventMap(func(agr countSizeAggregate) int64 {
		return agr.Count
	})

	sizePrint := fileExtPrint{
		count:   func(f *file) int64 { return e.total.extensions[f.path].Count },
		size:    func(f *file) uint64 { return uint64(f.size) },
		p:       p,
		headfmt: "\n<gray>TOP %d file extensions by size:</>\n\n",
		total:   e.total,
	}

	sizePrint.print(extBySize, e.top)

	countPrint := fileExtPrint{
		count:   func(f *file) int64 { return f.size },
		size:    func(f *file) uint64 { return e.total.extensions[f.path].Size },
		p:       p,
		headfmt: "\n<gray>TOP %d file extensions by count:</>\n\n",
		total:   e.total,
	}

	countPrint.print(extByCount, e.top)
}

type fileExtPrint struct {
	count   func(f *file) int64
	size    func(f *file) uint64
	p       printer
	headfmt string
	total   *totalInfo
}

func (fp *fileExtPrint) print(data files, top int) {
	const format = "%v\t%v\t%v\t%v\t%v\t%v\n"

	fp.p.cprint(fp.headfmt, top)

	fp.printTableHead(format)

	sort.Sort(sort.Reverse(data))

	for i := 0; i < top && i < len(data); i++ {
		h := data[i].path

		count := fp.count(data[i])
		sz := fp.size(data[i])

		fp.total.printCountAndSizeStatLine(fp.p, i+1, count, sz, h)
	}

	fp.p.flush()
}

func (fp *fileExtPrint) printTableHead(format string) {
	fp.p.tprint(format, " #", "Extension", "Count", "%", "Size", "%")
	fp.p.tprint(format, "--", "---------", "-----", "------", "----", "------")
}

func (e *extRenderer) evolventMap(mapper func(countSizeAggregate) int64) files {
	var result = make(files, len(e.total.extensions))
	i := 0
	for k, v := range e.total.extensions {
		result[i] = &file{size: mapper(v), path: k}
		i++
	}
	return result
}
