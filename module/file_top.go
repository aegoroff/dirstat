package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
)

func newTopFilesRenderer(work *topFilesWorker) renderer {
	return &topFilesRenderer{work}
}

type topFilesWorker struct {
	voidInit
	voidFinalize
	tree *fixedTree
}

type topFilesRenderer struct {
	*topFilesWorker
}

func newTopFilesWorker(top int) *topFilesWorker {
	return &topFilesWorker{tree: newFixedTree(top)}
}

// Worker methods

func (m *topFilesWorker) handler(evt *sys.ScanEvent) {
	f := evt.File

	fileContainer := file{size: f.Size, path: f.Path}
	m.tree.insert(&fileContainer)
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
