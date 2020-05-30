package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
)

// NewTopFilesModule creates new top files statistic module
func NewTopFilesModule(c *Context) Module {
	work := newTopFilesWorker(c.top)
	rend := &topFilesRenderer{*work}
	m := NewModule()
	m.addWorker(work)
	m.addRenderer(rend)
	return m
}

// NewTopFilesHiddenModule creates new top files statistic module
// that has disabled output
func NewTopFilesHiddenModule(c *Context) Module {
	work := newTopFilesWorker(c.top)
	m := NewModule()
	m.addWorker(work)
	return m
}

type topFilesWorker struct {
	tree *rbtree.RbTree
	top  int
}

type topFilesRenderer struct {
	topFilesWorker
}

func newTopFilesWorker(top int) *topFilesWorker {
	return &topFilesWorker{rbtree.NewRbTree(), top}
}

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

func (m *topFilesRenderer) print(p printer) {
	p.print("\nTOP %d files by size:\n\n", m.top)

	p.printtab("%v\t%v\n", "File", "Size")
	p.printtab("%v\t%v\n", "------", "----")

	i := 1

	m.tree.Descend(func(c rbtree.Comparable) bool {
		file := c.(*container)
		h := fmt.Sprintf("%d. %s", i, file.name)

		i++

		p.printtab("%v\t%v\n", h, human(file.size))

		return true
	})

	p.flush()
}
