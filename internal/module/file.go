package module

import (
	"github.com/aegoroff/godatastruct/rbtree"
)

type file struct {
	path string
	size int64
	pd   decorator
}

type files []*file

func (fi files) Len() int           { return len(fi) }
func (fi files) Less(i, j int) bool { return fi[i].size < fi[j].size }
func (fi files) Swap(i, j int)      { fi[i], fi[j] = fi[j], fi[i] }

func (f *file) Less(y rbtree.Comparable) bool  { return f.size < y.(*file).size }
func (f *file) Equal(y rbtree.Comparable) bool { return f.size == y.(*file).size }
func (f *file) String() string                 { return f.pd.decorate(f.path) }
