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
	p.cprint("\n<gray>The first file's size digit distribution of non zero files (benford law):</>\n\n")

	const format = "%v\t%v\t%v\t%v\t%v\t%v\n"
	p.tprint(format, "Digit", "Count", "%", "Benford ideal", "%", "Deviation")
	p.tprint(format, "-----", "-----", "------", "-------------", "---------", "---------")

	// IDEAL percents
	ideals := []float64{30.1, 17.6, 12.5, 9.7, 7.9, 6.7, 5.8, 5.1, 4.6}

	total := float64(b.total.FilesTotal.Count - b.distribution[0])
	for i, count := range b.distribution {
		if i == 0 {
			continue
		}
		percentOfCount := percent(float64(count), total)
		ideal := int64((ideals[i-1] / 100) * total)

		diff := count - ideal
		var deviation float64
		if ideal == 0 {
			deviation = 0
		} else {
			deviation = float64(diff) / float64(ideal)
		}

		p.tprint("%v\t%v\t%.2f%%\t%v\t%.2f%%\t%.2f%%\n", i, count, percentOfCount, ideal, ideals[i-1], deviation*100)
	}

	p.flush()
}
