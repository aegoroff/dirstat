package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"io"
	"text/tabwriter"
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

// NewFoldersModule creates new folders module
func NewFoldersModule(ctx *Context) Module {
	work := newFoldersWorker(ctx)
	rend := foldersRenderer{work}
	m := moduleFolders{
		work,
		rend,
	}
	return &m
}

// NewFoldersHiddenModule creates new folders module
// that has disabled output
func NewFoldersHiddenModule(ctx *Context) Module {
	work := newFoldersWorker(ctx)
	m := moduleFoldersNoOut{
		work,
		emptyRenderer{},
	}
	return &m
}

type foldersWorker struct {
	total    *totalInfo
	folders  *rbtree.RbTree
	topSize  *rbtree.RbTree
	topCount *rbtree.RbTree
}

type foldersRenderer struct {
	foldersWorker
}

type moduleFolders struct {
	foldersWorker
	foldersRenderer
}

type moduleFoldersNoOut struct {
	foldersWorker
	emptyRenderer
}

func newFoldersWorker(ctx *Context) foldersWorker {
	return foldersWorker{
		total:    ctx.total,
		folders:  rbtree.NewRbTree(),
		topSize:  rbtree.NewRbTree(),
		topCount: rbtree.NewRbTree(),
	}
}

func (m *foldersWorker) init() {
}

func (m *foldersWorker) postScan() {
	m.folders.WalkInorder(func(node *rbtree.Node) {
		fn := (*node.Key).(*folderNode)

		insertTo(m.topSize, &fn.container)

		fcn := folderCount{
			fn.container,
		}

		insertTo(m.topCount, &fcn)
	})

	m.total.CountFolders = m.folders.Root.Size
}

func (m *foldersWorker) folderHandler(fe *sys.FolderEntry) {
	var cmp rbtree.Comparable
	cmp = &folderNode{
		container{name: fe.Path, count: fe.Count, size: fe.Size},
	}
	m.folders.Insert(rbtree.NewNode(&cmp))
}

func (m *foldersWorker) fileHandler(_ *sys.FileEntry) {

}

func (f *foldersRenderer) output(tw *tabwriter.Writer, w io.Writer) {
	const format = "%v\t%v\t%v\t%v\t%v\n"

	_, _ = fmt.Fprintf(w, "\nTOP %d folders by size:\n\n", top)

	f.outputFoldersStatHeader(tw, format)

	i := 1

	f.topSize.Descend(func(c *rbtree.Comparable) bool {

		folder := (*c).(*container)
		f.outputFoldersStatLine(&i, folder, tw)

		return true
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d folders by count:\n\n", top)

	f.outputFoldersStatHeader(tw, format)

	i = 1

	f.topCount.Descend(func(c *rbtree.Comparable) bool {

		folder := (*c).(*folderCount)
		f.outputFoldersStatLine(&i, &folder.container, tw)

		return true
	})

	_ = tw.Flush()
}

func (f *foldersRenderer) outputFoldersStatLine(i *int, folder *container, tw *tabwriter.Writer) {
	h := fmt.Sprintf("%d. %s", *i, folder.name)

	*i++

	count := folder.count
	sz := uint64(folder.size)

	f.total.outputTopStatLine(tw, count, sz, h)
}

func (f *foldersRenderer) outputFoldersStatHeader(tw *tabwriter.Writer, format string) {
	_, _ = fmt.Fprintf(tw, format, "Folder", "Files", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "------", "-----", "------", "----", "------")
}
