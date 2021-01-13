package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/aegoroff/godatastruct/rbtree/special"
)

type extRenderer struct {
	*baseRenderer
	total   *totalInfo
	bySize  rbtree.RbTree
	byCount rbtree.RbTree
}

func newExtRenderer(ctx *Context, order int) renderer {
	return &extRenderer{
		baseRenderer: newBaseRenderer(order),
		total:        ctx.total,
		bySize:       special.NewMaxTree(int64(ctx.top)),
		byCount:      special.NewMaxTree(int64(ctx.top)),
	}
}

// Renderer method

func (e *extRenderer) render(p out.Printer) {
	e.createTops()

	heads := []string{"#", "Extension", "Count", "%", "Size", "%"}

	p.Cprint("\n<gray>TOP %d file extensions by size:</>\n\n", e.bySize.Len())

	ts := newTop(e.bySize, castSize, heads)
	ts.print(p, e.total)

	p.Cprint("\n<gray>TOP %d file extensions by count:</>\n\n", e.byCount.Len())

	tc := newTop(e.byCount, castCount, heads)
	tc.print(p, e.total)
}

func (e *extRenderer) createTops() {
	pd := &nonDestructiveDecorator{}
	for k, v := range e.total.extensions {
		fn := folder{
			path:  k,
			count: v.Count,
			size:  int64(v.Size),
			pd:    pd,
		}

		fs := folderS{fn}
		e.bySize.Insert(&fs)

		fc := folderC{fn}
		e.byCount.Insert(&fc)
	}
}
