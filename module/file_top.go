package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
)

func newTopFilesWorker(top int) *topFilesWorker {
	return &topFilesWorker{tree: newFixedTree(top)}
}

func newTopFilesRenderer(work *topFilesWorker) renderer {
	w := topFilesRenderer{work}

	w.fileFilterMiddleware = newFileFilterMiddleware(w.onFile)

	return &w
}

type topFilesWorker struct {
	voidInit
	voidFinalize
	*fileFilterMiddleware
	tree *fixedTree
}

type topFilesRenderer struct {
	*topFilesWorker
}

// Worker methods

func (m *topFilesWorker) onFile(f *sys.FileEntry) {
	fc := file{size: f.Size, path: f.Path}
	m.tree.insert(&fc)
}

// Renderer method

func (m *topFilesRenderer) print(p printer) {
	p.cprint("\n<gray>TOP %d files by size:</>\n\n", m.tree.size)

	p.print("%v\t%v\n", "File", "Size")
	p.print("%v\t%v\n", "------", "----")

	i := 1

	m.tree.tree.Descend(func(n rbtree.Node) bool {
		file := n.Key().(*file)
		h := fmt.Sprintf("%2d. %s", i, file)

		i++

		p.print("%v\t%v\n", h, human(file.size))

		return true
	})

	p.flush()
}
