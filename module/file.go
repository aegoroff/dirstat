package module

type file struct {
	path string
	size int64
}

type files []*file

func (fi files) Len() int           { return len(fi) }
func (fi files) Less(i, j int) bool { return fi[i].size < fi[j].size }
func (fi files) Swap(i, j int)      { fi[i], fi[j] = fi[j], fi[i] }

func (f *file) LessThan(y interface{}) bool { return f.size < y.(*file).size }
func (f *file) EqualTo(y interface{}) bool  { return f.size == y.(*file).size }
func (f *file) String() string              { return f.path }
