package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
)

type aggregateFile struct {
	ranges rbtree.RbTree
}

type aggregateFileHandler struct {
	*aggregateFile
}

type aggregateFileRenderer struct {
	*baseRenderer
	*aggregateFile
	total *totalInfo
}

func newAggregateFile(rs rbtree.RbTree) *aggregateFile {
	return &aggregateFile{
		ranges: rs,
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

// Handler method

func (m *aggregateFileHandler) Handle(evt *scan.Event) {
	f := evt.File

	n, ok := m.ranges.Floor(&Range{Min: f.Size})
	if ok {
		r := n.(*Range)
		r.size += f.Size
		r.count++
	}
}

// Renderer method

func (m *aggregateFileRenderer) render(p out.Printer) {
	p.Cprint("<gray>Total files stat:</>\n\n")

	topp := newTopper(p, m.total, []string{"#", "File size", "Amount", "%", "Size", "%"}, &noChangeDecorator{})
	topp.ascend(m.ranges)
}
