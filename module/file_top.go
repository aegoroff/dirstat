package module

import (
	"dirstat/module/internal/sys"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/aegoroff/godatastruct/rbtree/special"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func newTopFilesHandler(top int, pd decorator) *topFilesHandler {
	w := &topFilesHandler{
		tree: special.NewMaxTree(int64(top)),
		pd:   pd,
	}
	return w
}

func newTopFilesRenderer(work *topFilesHandler, order int) renderer {
	w := topFilesRenderer{
		topFilesHandler: work,
		baseRenderer:    newBaseRenderer(order),
	}

	return &w
}

type topFilesHandler struct {
	tree rbtree.RbTree
	pd   decorator
}

type topFilesRenderer struct {
	*topFilesHandler
	*baseRenderer
}

// Worker method

func (m *topFilesHandler) Handle(evt *sys.ScanEvent) {
	f := evt.File
	fc := file{size: f.Size, path: f.Path, pd: m.pd}
	m.tree.Insert(&fc)
}

// Renderer method

func (m *topFilesRenderer) print(p printer) {
	p.cprint("\n<gray>TOP %d files by size:</>\n\n", m.tree.Len())

	tab := p.createTab()

	tab.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignLeft, WidthMax: 100},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: sizeTransformer},
	})

	appendHeaders([]string{"#", "File", "Size"}, tab)

	i := 1

	it := rbtree.NewDescend(m.tree).Iterator()

	for it.Next() {
		n := it.Current()
		file, ok := n.(*file)

		if !ok {
			p.cprint("<red>Invalid casting: expected *file key type but it wasn`t</>")
			return
		}

		tab.AppendRow([]interface{}{
			i,
			file,
			uint64(file.size),
		})

		i++
	}

	tab.Render()
}
