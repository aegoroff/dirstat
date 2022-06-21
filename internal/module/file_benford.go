package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type benfordFile struct {
	distribution []int64
}

type benfordFileHandler struct {
	*benfordFile
}

type benfordFileRenderer struct {
	*baseRenderer
	*benfordFile
	total *totalInfo
}

func newBenfordFile() *benfordFile {
	return &benfordFile{
		distribution: []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func newBenfordFileHandler(bf *benfordFile) scan.Handler {
	return newOnlyFilesHandler(&benfordFileHandler{bf})
}

func newBenfordFileRenderer(ctx *Context, bf *benfordFile, order int) renderer {
	return &benfordFileRenderer{
		baseRenderer: newBaseRenderer(order),
		benfordFile:  bf,
		total:        ctx.total,
	}
}

// Handler method

func (b *benfordFileHandler) Handle(evt *scan.Event) {
	s := evt.File.Size
	for s > 9 {
		s = div10(s)
	}
	b.distribution[s]++
}

func div10(n int64) int64 {
	var q int64
	var r int64
	q = (n >> 1) + (n >> 2)
	q = q + (q >> 4)
	q = q + (q >> 8)
	q = q + (q >> 16)
	q = q >> 3
	r = n - (((q << 2) + q) << 1)
	if r > 9 {
		return q + 1
	}
	return q
}

// Renderer method

func (b *benfordFileRenderer) render(p out.Printer) {
	p.Println()
	p.Cprint("<gray>The first file's size digit distribution of non zero files (benford law):</>")
	p.Println()
	p.Println()

	tw := newTableWriter(p)

	tw.configColumns([]table.ColumnConfig{
		{Number: 1, Align: text.AlignLeft, AlignHeader: text.AlignLeft},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.countTransformer},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
		{Number: 4, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.countTransformer},
		{Number: 5, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
		{Number: 6, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
	})

	tw.addHeaders([]string{"Digit", "Count", "%", "Benford ideal", "%", "Deviation"})

	// IDEAL percents
	ideals := []float64{30.1, 17.6, 12.5, 9.7, 7.9, 6.7, 5.8, 5.1, 4.6}

	total := float64(b.total.FilesTotal - b.distribution[0])
	for i, count := range b.distribution {
		if i == 0 {
			continue
		}

		percentOfCount := percent(float64(count), total)
		ideal := int64((ideals[i-1] / 100) * total)

		diff := count - ideal
		var deviation float64
		if diff == 0 {
			deviation = 0
		} else if ideal == 0 {
			deviation = 1
		} else {
			deviation = float64(diff) / float64(ideal)
		}

		tw.addRow([]interface{}{
			i,
			count,
			percentOfCount,
			ideal,
			ideals[i-1],
			deviation * 100,
		})
	}

	tw.render()
}
