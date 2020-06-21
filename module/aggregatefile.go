package module

import "dirstat/module/internal/sys"

type aggregateFileWorker struct {
	aggregate map[Range]fileStat
}

type aggregateFileRenderer struct {
	work  *aggregateFileWorker
	total *totalInfo
}

func newAggregateFileWorker() *aggregateFileWorker {
	return &aggregateFileWorker{
		aggregate: make(map[Range]fileStat),
	}
}

func newAggregateFileRenderer(ctx *Context, w *aggregateFileWorker) *aggregateFileRenderer {
	return &aggregateFileRenderer{
		total: ctx.total,
		work:  w,
	}
}

func (*aggregateFileWorker) init()     {}
func (*aggregateFileWorker) finalize() {}

func (m *aggregateFileWorker) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}
	f := evt.File

	unsignedSize := uint64(f.Size)

	// Calculate files range statistic
	for _, r := range fileSizeRanges {
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

	heads := createRangesHeads()
	for i, r := range fileSizeRanges {
		count := m.work.aggregate[r].TotalFilesCount
		sz := m.work.aggregate[r].TotalFilesSize

		m.total.printCountAndSizeStatLine(p, count, sz, heads[i])
	}
	p.flush()
}
