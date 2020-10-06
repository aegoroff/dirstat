package module

import (
	"dirstat/module/internal/sys"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func newTopFilesWorker(top int, pd *pathDecorator) *topFilesWorker {
	return &topFilesWorker{tree: newFixedTree(top), pd: pd}
}

func newTopFilesRenderer(work *topFilesWorker) renderer {
	w := topFilesRenderer{topFilesWorker: work}

	w.fileFilter = newFileFilter(w.onFile)

	return &w
}

type topFilesWorker struct {
	voidInit
	voidFinalize
	*fileFilter
	tree *fixedTree
	pd   *pathDecorator
}

type topFilesRenderer struct {
	*topFilesWorker
}

// Worker methods

func (m *topFilesWorker) onFile(f *sys.FileEntry) {
	fc := file{size: f.Size, path: f.Path, pd: m.pd}
	m.tree.insert(&fc)
}

// Renderer method

func (m *topFilesRenderer) print(p printer) {
	p.cprint("\n<gray>TOP %d files by size:</>\n\n", m.tree.size)

	tab := p.createTab()

	tab.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignLeft, WidthMax: 100},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: sizeTransformer},
	})

	appendHeaders([]string{"#", "File", "Size"}, tab)

	i := 1

	m.tree.tree.Descend(func(n rbtree.Node) bool {
		file, ok := n.Key().(*file)

		if !ok {
			p.cprint("<red>Invalid casting: expected *file key type but it wasn`t</>")
			return false
		}

		tab.AppendRow([]interface{}{
			i,
			file,
			uint64(file.size),
		})

		i++
		return true
	})

	tab.Render()
}
