package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/aegoroff/godatastruct/rbtree/special"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type topFiles struct {
	tree rbtree.RbTree
}

type topFilesHandler struct {
	*topFiles
}

type topFilesRenderer struct {
	*topFiles
	*baseRenderer
	pd decorator
}

func newTopFiles(top int) *topFiles {
	return &topFiles{
		tree: special.NewMaxTree(int64(top)),
	}
}

func newTopFilesHandler(tf *topFiles) scan.Handler {
	h := &topFilesHandler{
		topFiles: tf,
	}
	return newOnlyFilesHandler(h)
}

func newTopFilesRenderer(tf *topFiles, pd decorator, order int) renderer {
	w := topFilesRenderer{
		topFiles:     tf,
		baseRenderer: newBaseRenderer(order),
		pd:           pd,
	}

	return &w
}

// Handler method

func (m *topFilesHandler) Handle(evt *scan.Event) {
	f := evt.File
	fc := file{size: f.Size, path: f.Path}
	m.tree.Insert(&fc)
}

// Renderer method

func (m *topFilesRenderer) render(p out.Printer) {
	p.Cprint("\n<gray>TOP %d files by size:</>\n\n", m.tree.Len())

	tw := newTableWriter(p)

	tw.configColumns([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: m.decoratePathOrName, WidthMax: 100},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.sizeTransformer},
	})

	tw.addHeaders([]string{"#", "File", "Size"})

	i := 1

	it := rbtree.NewDescend(m.tree).Iterator()

	for it.Next() {
		n := it.Current()
		file, ok := n.(*file)

		if !ok {
			p.Cprint("<red>Invalid casting: expected *file key type but it wasn`t</>")
			return
		}

		tw.addRow([]any{
			i,
			file.String(),
			file.size,
		})

		i++
	}

	tw.render()
}

func (m *topFilesRenderer) decoratePathOrName(val any) string {
	return m.pd.decorate(val.(string))
}
