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
	e.fillTops()

	heads := []string{"#", "Extension", "Count", "%", "Size", "%"}
	top := newTopper(p, e.total, heads, &noChangeDecorator{})

	p.Cprint("\n<gray>TOP %d file extensions by size:</>\n\n", e.bySize.Len())
	top.descend(e.bySize)

	p.Cprint("\n<gray>TOP %d file extensions by count:</>\n\n", e.byCount.Len())
	top.descend(e.byCount)
}

func (e *extRenderer) fillTops() {
	it := rbtree.NewWalkInorder(e.total.extensions)
	it.Foreach(func(k rbtree.Comparable) {
		fn := k.(*folder)
		fs := folderS{fn}
		e.bySize.Insert(&fs)

		fc := folderC{fn}
		e.byCount.Insert(&fc)
	})
}
