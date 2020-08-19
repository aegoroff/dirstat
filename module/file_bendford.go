package module

import (
	"dirstat/module/internal/sys"
)

type bendfordFileWorker struct {
	voidFinalize
	voidInit
	*fileFilter
	distribution []int64
	total        *totalInfo
}

type bendfordFileRenderer struct {
	*bendfordFileWorker
}

func newBendfordFileWorker(ctx *Context) *bendfordFileWorker {
	w := bendfordFileWorker{
		distribution: []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		total:        ctx.total,
	}
	w.fileFilter = newFileFilter(w.onFile)
	return &w
}

func newBendfordFileRenderer(work *bendfordFileWorker) renderer {
	return &bendfordFileRenderer{work}
}

func (b *bendfordFileWorker) onFile(f *sys.FileEntry) {
	s := f.Size
	for s >= 10 {
		s = s / 10
	}
	b.distribution[s]++
}

func (b *bendfordFileRenderer) print(p printer) {
	p.cprint("\n<gray>File sizes by first file's size digit distribution (bendford law):</>\n\n")

	const format = "%v\t%v\t%v\n"
	p.tprint(format, "Digit", "Count", "%")
	p.tprint(format, "-----", "-----", "------")

	for i, count := range b.distribution {
		percentOfCount := b.total.countPercent(count)
		p.tprint("%v\t%v\t%.2f%%\n", i, count, percentOfCount)
	}

	p.flush()
}
