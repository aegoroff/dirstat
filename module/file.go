package module

import "dirstat/module/internal/sys"

type file struct {
	path string
	size int64
	pd   *pathDecorator
}

type files []*file

type fileHandler func(f *sys.FileEntry)

type fileFilter struct {
	h fileHandler
}

func (fi files) Len() int           { return len(fi) }
func (fi files) Less(i, j int) bool { return fi[i].size < fi[j].size }
func (fi files) Swap(i, j int)      { fi[i], fi[j] = fi[j], fi[i] }

func (f *file) LessThan(y interface{}) bool { return f.size < y.(*file).size }
func (f *file) EqualTo(y interface{}) bool  { return f.size == y.(*file).size }
func (f *file) String() string {
	if f.pd == nil {
		return f.path
	}
	return f.pd.decorate(f.path)
}

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
