package module

import (
	"github.com/aegoroff/godatastruct/rbtree"
)

// folder represents file system container that described by path
// and has size and the number of elements in it (count field).
type folder struct {
	path  string
	size  int64
	count int64
}

// Count sortable folder
type folderC struct {
	folder
}

// Size sortable folder
type folderS struct {
	folder
}

// Path sortable folder methods

func (f *folder) String() string { return f.path }
func (f *folder) Size() int64    { return f.size }
func (f *folder) Count() int64   { return f.count }

// Count sortable folder methods

func (fc *folderC) Less(y rbtree.Comparable) bool  { return fc.count < y.(*folderC).count }
func (fc *folderC) Equal(y rbtree.Comparable) bool { return fc.count == y.(*folderC).count }

// Size sortable folder methods

func (fs *folderS) Less(y rbtree.Comparable) bool  { return fs.size < y.(*folderS).size }
func (fs *folderS) Equal(y rbtree.Comparable) bool { return fs.size == y.(*folderS).size }
