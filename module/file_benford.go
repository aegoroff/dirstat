package module

import (
	"dirstat/module/internal/sys"
)

type benfordFileWorker struct {
	voidFinalize
	voidInit
	*fileFilter
	distribution []int64
	total        *totalInfo
}

type benfordFileRenderer struct {
	*benfordFileWorker
}

func newBenfordFileWorker(ctx *Context) *benfordFileWorker {
	w := benfordFileWorker{
		distribution: []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		total:        ctx.total,
	}
	w.fileFilter = newFileFilter(w.onFile)
	return &w
}

func newBenfordFileRenderer(work *benfordFileWorker) renderer {
	return &benfordFileRenderer{work}
}

func (b *benfordFileWorker) onFile(f *sys.FileEntry) {
	s := f.Size
	for s >= 10 {
		s = s / 10
	}
	b.distribution[s]++
}

func (b *benfordFileRenderer) print(p printer) {
	p.cprint("\n<gray>File sizes by first file's size digit distribution (benford law):</>\n\n")

	const format = "%v\t%v\t%v\n"
	p.tprint(format, "Digit", "Count", "%")
	p.tprint(format, "-----", "-----", "------")

	for i, count := range b.distribution {
		percentOfCount := b.total.countPercent(count)
		p.tprint("%v\t%v\t%.2f%%\n", i, count, percentOfCount)
	}

	p.flush()
}
