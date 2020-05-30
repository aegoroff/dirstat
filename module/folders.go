package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
)

type folderNode struct {
	container
}

type folderCount struct {
	container
}

func (f *folderNode) LessThan(y interface{}) bool {
	return f.name < y.(*folderNode).name
}

func (f *folderNode) EqualTo(y interface{}) bool {
	return f.name == y.(*folderNode).name
}

func (f *folderCount) LessThan(y interface{}) bool {
	return f.count < y.(*folderCount).count
}

func (f *folderCount) EqualTo(y interface{}) bool {
	return f.count == y.(*folderCount).count
}

type foldersWorker struct {
	total   *totalInfo
	folders rbtree.RbTree
	bySize  rbtree.RbTree
	byCount rbtree.RbTree
	top     int
}

type foldersRenderer struct {
	work *foldersWorker
}

func newFoldersWorker(ctx *Context) *foldersWorker {
	return &foldersWorker{
		total:   ctx.total,
		folders: rbtree.NewRbTree(),
		bySize:  rbtree.NewRbTree(),
		byCount: rbtree.NewRbTree(),
		top:     ctx.top,
	}
}

func newFoldersRenderer(work *foldersWorker) renderer {
	return &foldersRenderer{work}
}

// Worker methods

func (m *foldersWorker) init() {
}

func (m *foldersWorker) finalize() {
	m.folders.WalkInorder(func(node rbtree.Comparable) {
		fn := node.(*folderNode)

		insertTo(m.bySize, m.top, &fn.container)

		fcn := folderCount{
			fn.container,
		}

		insertTo(m.byCount, m.top, &fcn)
	})

	m.total.CountFolders = m.folders.Len()
}

func (m *foldersWorker) handler(evt *sys.ScanEvent) {
	if evt.Folder == nil {
		return
	}
	fe := evt.Folder

	fn := folderNode{
		container{name: fe.Path, count: fe.Count, size: fe.Size},
	}
	m.folders.Insert(&fn)
}

// Renderer method

func (f *foldersRenderer) print(p printer) {
	const format = "%v\t%v\t%v\t%v\t%v\n"

	p.print("\nTOP %d folders by size:\n\n", f.work.top)

	f.printTableHead(p, format)

	i := 1

	f.work.bySize.Descend(func(c rbtree.Comparable) bool {

		folder := c.(*container)
		f.printTableRow(&i, folder, p)

		return true
	})

	p.flush()

	p.print("\nTOP %d folders by count:\n\n", f.work.top)

	f.printTableHead(p, format)

	i = 1

	f.work.byCount.Descend(func(c rbtree.Comparable) bool {

		folder := c.(*folderCount)
		f.printTableRow(&i, &folder.container, p)

		return true
	})

	p.flush()
}

func (f *foldersRenderer) printTableRow(i *int, folder *container, p printer) {
	h := fmt.Sprintf("%d. %s", *i, folder.name)

	*i++

	count := folder.count
	sz := uint64(folder.size)

	f.work.total.printCountAndSizeStatLine(p, count, sz, h)
}

func (f *foldersRenderer) printTableHead(rc printer, format string) {
	rc.printtab(format, "Folder", "Files", "%", "Size", "%")
	rc.printtab(format, "------", "-----", "------", "----", "------")
}
