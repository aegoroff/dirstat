package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/aegoroff/godatastruct/rbtree/special"
)

type foldersHandler struct {
	*folders
}

type foldersRenderer struct {
	*folders
	*baseRenderer
	total *totalInfo
	pd    decorator
}

type folders struct {
	bySize  rbtree.RbTree
	byCount rbtree.RbTree
}

func newFolders(top int) *folders {
	return &folders{
		bySize:  special.NewMaxTree(int64(top)),
		byCount: special.NewMaxTree(int64(top)),
	}
}

func newFoldersHandler(fc *folders) scan.Handler {
	h := &foldersHandler{fc}
	return newOnlyFoldersHandler(h)
}

func newFoldersRenderer(f *folders, ctx *Context, order int) renderer {
	return &foldersRenderer{
		folders:      f,
		total:        ctx.total,
		baseRenderer: newBaseRenderer(order),
		pd:           ctx.pd,
	}
}

// Handler method

func (m *foldersHandler) Handle(evt *scan.Event) {
	fe := evt.Folder

	fn := &folder{
		path:  fe.Path,
		count: fe.Count,
		size:  fe.Size,
	}

	fs := folderS{fn}
	m.bySize.Insert(&fs)

	fc := folderC{fn}
	m.byCount.Insert(&fc)
}

// Renderer method

func (f *foldersRenderer) render(p out.Printer) {
	heads := []string{"#", "Folder", "Files", "%", "Size", "%"}
	top := newTopper(p, f.total, heads, f.pd)

	p.Cprint("\n<gray>TOP %d folders by size:</>\n\n", f.bySize.Len())
	top.descend(f.bySize)

	p.Cprint("\n<gray>TOP %d folders by count:</>\n\n", f.byCount.Len())
	top.descend(f.byCount)
}
