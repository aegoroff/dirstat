package module

import "dirstat/module/internal/sys"

type aggregateFileWorker struct {
	voidInit
	voidFinalize
	*fileFilter
	aggregate  map[Range]fileStat
	fileRanges ranges
}

type aggregateFileRenderer struct {
	work  *aggregateFileWorker
	total *totalInfo
}

func newAggregateFileWorker(rs ranges) *aggregateFileWorker {
	w := aggregateFileWorker{
		aggregate:  make(map[Range]fileStat, len(rs)),
		fileRanges: rs,
	}

	w.fileFilter = newFileFilter(w.onFile)

	return &w
}

func newAggregateFileRenderer(ctx *Context, w *aggregateFileWorker) *aggregateFileRenderer {
	return &aggregateFileRenderer{
		total: ctx.total,
		work:  w,
	}
}

func (m *aggregateFileWorker) onFile(f *sys.FileEntry) {
	unsignedSize := uint64(f.Size)

	// Calculate files range statistic
	for _, r := range m.fileRanges {
		if !r.Contains(f.Size) {
			continue
		}

		s := m.aggregate[r]
		s.TotalFilesCount++
		s.TotalFilesSize += unsignedSize
		m.aggregate[r] = s
	}
}

// Renderer method

func (m *aggregateFileRenderer) print(p printer) {
	const format = "%v\t%v\t%v\t%v\t%v\n"

	p.cprint("<gray>Total files stat:</>\n\n")
	p.print(format, "File size", "Amount", "%", "Size", "%")
	p.print(format, "---------", "------", "------", "----", "------")

	heads := m.work.fileRanges.heads()
	for i, r := range m.work.fileRanges {
		count := m.work.aggregate[r].TotalFilesCount
		sz := m.work.aggregate[r].TotalFilesSize

		m.total.printCountAndSizeStatLine(p, count, sz, heads[i])
	}
	p.flush()
}
