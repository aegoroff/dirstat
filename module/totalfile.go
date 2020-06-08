package module

import "dirstat/module/internal/sys"

type totalFileWorker struct {
	aggregate map[Range]fileStat
}

type totalFileRenderer struct {
	work  *totalFileWorker
	total *totalInfo
}

func newTotalFileWorker() *totalFileWorker {
	return &totalFileWorker{
		aggregate: make(map[Range]fileStat),
	}
}

func newTotalFileRenderer(ctx *Context, w *totalFileWorker) *totalFileRenderer {
	return &totalFileRenderer{
		total: ctx.total,
		work:  w,
	}
}

func (*totalFileWorker) init()     {}
func (*totalFileWorker) finalize() {}

func (m *totalFileWorker) handler(evt *sys.ScanEvent) {
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

func (m *totalFileRenderer) print(p printer) {
	p.print("Total files stat:\n\n")

	const format = "%v\t%v\t%v\t%v\t%v\n"

	p.printtab(format, "File size", "Amount", "%", "Size", "%")
	p.printtab(format, "---------", "------", "------", "----", "------")

	heads := createRangesHeads()
	for i, r := range fileSizeRanges {
		count := m.work.aggregate[r].TotalFilesCount
		sz := m.work.aggregate[r].TotalFilesSize

		m.total.printCountAndSizeStatLine(p, count, sz, heads[i])
	}
	p.flush()
}
