package module

import (
	"dirstat/module/internal/sys"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/aegoroff/godatastruct/rbtree/special"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func newTopFilesWorker(top int, pd decorator) *topFilesWorker {
	return &topFilesWorker{tree: special.NewMaxTree(int64(top)), pd: pd}
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
	tree rbtree.RbTree
	pd   decorator
}

type topFilesRenderer struct {
	*topFilesWorker
}

// Worker methods

func (m *topFilesWorker) onFile(f *sys.FileEntry) {
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

	rbtree.NewDescend(m.tree).Foreach(func(n rbtree.Node) {
		file, ok := n.Key().(*file)

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
	})

	tab.Render()
}
