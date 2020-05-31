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
	tree rbtree.RbTree
	top  int
}

type topFilesRenderer struct {
	work *topFilesWorker
}

func newTopFilesWorker(top int) *topFilesWorker {
	return &topFilesWorker{rbtree.NewRbTree(), top}
}

// Worker methods

func (m *topFilesWorker) init() {
}

func (m *topFilesWorker) finalize() {

}

func (m *topFilesWorker) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}
	f := evt.File

	fileContainer := container{size: f.Size, name: f.Path, count: 1}
	insertTo(m.tree, m.top, &fileContainer)
}

// Renderer method

func (m *topFilesRenderer) print(p printer) {
	p.print("\nTOP %d files by size:\n\n", m.work.top)

	p.printtab("%v\t%v\n", "File", "Size")
	p.printtab("%v\t%v\n", "------", "----")

	i := 1

	m.work.tree.Descend(func(c rbtree.Node) bool {
		file := c.(rbtree.Comparable).(*container)
		h := fmt.Sprintf("%d. %s", i, file.name)

		i++

		p.printtab("%v\t%v\n", h, human(file.size))

		return true
	})

	p.flush()
}
