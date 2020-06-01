package module

import "github.com/aegoroff/godatastruct/rbtree"

type file struct {
	path string
	size int64
}

type files []*file

func (fi files) Len() int {
	return len(fi)
}

func (fi files) Less(i, j int) bool {
	return fi[i].size < fi[j].size
}

func (fi files) Swap(i, j int) {
	fi[i], fi[j] = fi[j], fi[i]
}

func (f *file) LessThan(y interface{}) bool {
	return f.size < y.(*file).size
}

func (f *file) EqualTo(y interface{}) bool {
	return f.size == y.(*file).size
}

func (f *file) String() string {
	return f.path
}

// insertTo inserts node into tree which size is limited by the size parameter.
// Only <size> max nodes will be in the tree
func insertTo(tree rbtree.RbTree, size int, c rbtree.Comparable) {
	min := tree.Minimum()
	if tree.Len() < int64(size) || min.Key().LessThan(c) {
		if tree.Len() == int64(size) {
			tree.DeleteNode(min.Key())
		}

		tree.Insert(c)
	}
}
