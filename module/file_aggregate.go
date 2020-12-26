package module

import (
	"dirstat/module/internal/sys"
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

func newAggregateFileHandler(af *aggregateFile) sys.Handler {
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

func (m *aggregateFileHandler) Handle(evt *sys.ScanEvent) {
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

func (m *aggregateFileRenderer) print(p printer) {
	p.cprint("<gray>Total files stat:</>\n\n")

	tab := p.createTab()

	tab.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignLeft, WidthMax: 100},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignLeft},
		{Number: 4, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: percentTransformer},
		{Number: 5, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: sizeTransformer},
		{Number: 6, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: percentTransformer},
	})

	appendHeaders([]string{"#", "File size", "Amount", "%", "Size", "%"}, tab)

	heads := m.fileRanges.heads(transparentDecorator)
	for i, r := range m.fileRanges {
		count := m.aggregate[r].TotalFilesCount
		sz := m.aggregate[r].TotalFilesSize

		percentOfCount := m.total.countPercent(count)
		percentOfSize := m.total.sizePercent(sz)

		tab.AppendRow([]interface{}{
			i + 1,
			heads[i],
			count,
			percentOfCount,
			sz,
			percentOfSize,
		})
	}
	tab.Render()
}
