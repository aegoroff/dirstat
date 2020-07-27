package module

import (
	"dirstat/module/internal/sys"
	"github.com/aegoroff/godatastruct/rbtree"
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

	const format = "%v\t%v\t%v\n"

	p.tprint(format, " #", "File", "Size")
	p.tprint(format, "--", "------", "----")

	i := 1

	m.tree.tree.Descend(func(n rbtree.Node) bool {
		file := n.Key().(*file)

		p.tprint("%2d\t%v\t%v\n", i, file, human(file.size))
		i++
		return true
	})

	p.flush()
}
