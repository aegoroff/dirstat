package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type aggregateFile struct {
	aggregate  map[Range]fileStat
	fileRanges ranges
}

type aggregateFileHandler struct {
	*aggregateFile
}

type aggregateFileRenderer struct {
	*baseRenderer
	*aggregateFile
	total *totalInfo
}

func newAggregateFile(rs ranges) *aggregateFile {
	return &aggregateFile{
		aggregate:  make(map[Range]fileStat, len(rs)),
		fileRanges: rs,
	}
}

func newAggregateFileHandler(af *aggregateFile) scan.Handler {
	return newOnlyFilesHandler(&aggregateFileHandler{af})
}

func newAggregateFileRenderer(ctx *Context, af *aggregateFile, order int) *aggregateFileRenderer {
	return &aggregateFileRenderer{
		baseRenderer:  newBaseRenderer(order),
		total:         ctx.total,
		aggregateFile: af,
	}
}

// Worker method

func (m *aggregateFileHandler) Handle(evt *scan.Event) {
	f := evt.File
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

func (m *aggregateFileRenderer) render(p out.Printer) {
	p.Cprint("<gray>Total files stat:</>\n\n")

	tw := newTableWriter(p)

	tw.configColumns([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignLeft, WidthMax: 100},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignLeft},
		{Number: 4, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
		{Number: 5, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.sizeTransformer},
		{Number: 6, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
	})

	tw.addHeaders([]string{"#", "File size", "Amount", "%", "Size", "%"})

	heads := m.fileRanges.heads(transparentDecorator)
	for i, r := range m.fileRanges {
		count := m.aggregate[r].TotalFilesCount
		sz := m.aggregate[r].TotalFilesSize

		percentOfCount := m.total.countPercent(count)
		percentOfSize := m.total.sizePercent(sz)

		tw.addRow([]interface{}{
			i + 1,
			heads[i],
			count,
			percentOfCount,
			sz,
			percentOfSize,
		})
	}
	tw.render()
}
