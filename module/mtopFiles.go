package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
)

// NewTopFilesModule creates new top files statistic module
func NewTopFilesModule(_ *Context) Module {
	work := newTopFilesWorker()
	rend := topFilesRenderer{work}
	m := moduleTopFiles{
		work,
		rend,
	}
	return &m
}

// NewTopFilesHiddenModule creates new top files statistic module
// that has disabled output
func NewTopFilesHiddenModule(_ *Context) Module {
	work := newTopFilesWorker()
	m := moduleTopFilesNoOut{
		work,
		emptyRenderer{},
	}
	return &m
}

type topFilesWorker struct {
	tree *rbtree.RbTree
}

type topFilesRenderer struct {
	topFilesWorker
}

type moduleTopFiles struct {
	topFilesWorker
	topFilesRenderer
}

type moduleTopFilesNoOut struct {
	topFilesWorker
	emptyRenderer
}

func newTopFilesWorker() topFilesWorker {
	return topFilesWorker{rbtree.NewRbTree()}
}

func (m *topFilesWorker) init() {
}

func (m *topFilesWorker) postScan() {

}

func (m *topFilesWorker) folderHandler(*sys.FolderEntry) {

}

func (m *topFilesWorker) fileHandler(f *sys.FileEntry) {
	fileContainer := container{size: f.Size, name: f.Path, count: 1}
	insertTo(m.tree, &fileContainer)
}

func (m *topFilesRenderer) output(ctx renderContext) {
	ctx.write("\nTOP %d files by size:\n\n", top)

	ctx.writetab("%v\t%v\n", "File", "Size")
	ctx.writetab("%v\t%v\n", "------", "----")

	i := 1

	m.tree.Descend(func(c rbtree.Comparable) bool {
		file := c.(*container)
		h := fmt.Sprintf("%d. %s", i, file.name)

		i++

		sz := uint64(file.size)

		ctx.writetab("%v\t%v\n", h, humanize.IBytes(sz))

		return true
	})

	ctx.flush()
}
