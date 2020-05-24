package cmd

// container represents file system container that described by name
// and has size and the number of elements in it (count field). It the case of file
// the number is 1 and if it's a folder count will be the number of files in it
type container struct {
	name  string
	size  int64
	count int64
}

type containers []*container

func (x containers) Len() int {
	return len(x)
}

func (x containers) Less(i, j int) bool {
	return x[i].size < x[j].size
}

func (x containers) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x *container) LessThan(y interface{}) bool {
	return x.size < (y.(*container)).size
}

func (x *container) EqualTo(y interface{}) bool {
	return x.size == (y.(*container)).size
}
