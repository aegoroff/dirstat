package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
)

// Folder interface
type folderI interface {
	Path() string
	Size() int64
	Count() int64
}

// folder represents file system container that described by path
// and has size and the number of elements in it (count field).
type folder struct {
	path  string
	size  int64
	count int64
}

// Count sortable folder
type folderC struct {
	path  string
	size  int64
	count int64
}

// Size sortable folder
type folderS struct {
	path  string
	size  int64
	count int64
}

// Path sortable folder methods

func (f *folder) LessThan(y interface{}) bool {
	return f.String() < y.(*folder).String()
}

func (f *folder) EqualTo(y interface{}) bool {
	return f.String() == y.(*folder).String()
}

func (f *folder) String() string {
	return f.path
}

func (f *folder) Path() string {
	return f.path
}

func (f *folder) Size() int64 {
	return f.size
}

func (f *folder) Count() int64 {
	return f.count
}

// Count sortable folder methods

func (fc *folderC) LessThan(y interface{}) bool {
	return fc.count < y.(*folderC).count
}

func (fc *folderC) EqualTo(y interface{}) bool {
	return fc.count == y.(*folderC).count
}

func (fc *folderC) String() string {
	return fc.path
}

func (fc *folderC) Path() string {
	return fc.path
}

func (fc *folderC) Size() int64 {
	return fc.size
}

func (fc *folderC) Count() int64 {
	return fc.count
}

// Size sortable folder methods

func (fs *folderS) LessThan(y interface{}) bool {
	return fs.size < y.(*folderS).size
}

func (fs *folderS) EqualTo(y interface{}) bool {
	return fs.size == y.(*folderS).size
}

func (fs *folderS) String() string {
	return fs.path
}

func (fs *folderS) Path() string {
	return fs.path
}

func (fs *folderS) Size() int64 {
	return fs.size
}

func (fs *folderS) Count() int64 {
	return fs.count
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
	m.folders.WalkInorder(func(node rbtree.Node) {
		fn := node.Key().(*folder)

		fs := folderS{
			path:  fn.path,
			count: fn.count,
			size:  fn.size,
		}
		insertTo(m.bySize, m.top, &fs)

		fc := folderC{
			path:  fn.path,
			count: fn.count,
			size:  fn.size,
		}

		insertTo(m.byCount, m.top, &fc)
	})

	m.total.CountFolders = m.folders.Len()
}

func (m *foldersWorker) handler(evt *sys.ScanEvent) {
	if evt.Folder == nil {
		return
	}
	fe := evt.Folder

	fn := folder{
		path:  fe.Path,
		count: fe.Count,
		size:  fe.Size,
	}
	m.folders.Insert(&fn)
}

// Renderer method

func (f *foldersRenderer) print(p printer) {
	const format = "%v\t%v\t%v\t%v\t%v\n"

	p.print("\nTOP %d folders by size:\n\n", f.work.top)

	f.printTableHead(p, format)

	i := 1

	f.work.bySize.Descend(func(n rbtree.Node) bool {

		folder := n.Key().(*folderS)
		f.printTableRow(&i, folder, p)

		return true
	})

	p.flush()

	p.print("\nTOP %d folders by count:\n\n", f.work.top)

	f.printTableHead(p, format)

	i = 1

	f.work.byCount.Descend(func(n rbtree.Node) bool {

		folder := n.Key().(*folderC)
		f.printTableRow(&i, folder, p)

		return true
	})

	p.flush()
}

func (f *foldersRenderer) printTableRow(i *int, fi folderI, p printer) {
	h := fmt.Sprintf("%d. %s", *i, fi.Path())

	*i++

	count := fi.Count()
	sz := uint64(fi.Size())

	f.work.total.printCountAndSizeStatLine(p, count, sz, h)
}

func (f *foldersRenderer) printTableHead(rc printer, format string) {
	rc.printtab(format, "Folder", "Files", "%", "Size", "%")
	rc.printtab(format, "------", "-----", "------", "----", "------")
}
