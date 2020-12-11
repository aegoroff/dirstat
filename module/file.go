package module

import (
	"dirstat/module/internal/sys"
	"github.com/aegoroff/godatastruct/rbtree"
)

type file struct {
	path string
	size int64
	pd   decorator
}

type files []*file

type fileHandler func(f *sys.FileEntry)

type fileFilter struct {
	h fileHandler
}

func (fi files) Len() int           { return len(fi) }
func (fi files) Less(i, j int) bool { return fi[i].size < fi[j].size }
func (fi files) Swap(i, j int)      { fi[i], fi[j] = fi[j], fi[i] }

func (f *file) LessThan(y rbtree.Comparable) bool { return f.size < y.(*file).size }
func (f *file) EqualTo(y rbtree.Comparable) bool  { return f.size == y.(*file).size }
func (f *file) String() string                    { return f.pd.decorate(f.path) }

func newFileFilter(h fileHandler) *fileFilter {
	return &fileFilter{
		h: h,
	}
}

func (f *fileFilter) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}
	f.h(evt.File)
}
