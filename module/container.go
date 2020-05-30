package module

import "github.com/aegoroff/godatastruct/rbtree"

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

func (c *container) LessThan(y interface{}) bool {
	return c.size < (y.(*container)).size
}

func (c *container) EqualTo(y interface{}) bool {
	return c.size == (y.(*container)).size
}

// insertTo inserts node into tree which size is limited by the size parameter.
// Only <size> max nodes will be in the tree
func insertTo(tree rbtree.RbTree, size int, c rbtree.Comparable) {
	min := tree.Minimum()
	if tree.Len() < int64(size) || min.LessThan(c) {
		if tree.Len() == int64(size) {
			tree.DeleteNode(min)
		}

		tree.Insert(c)
	}
}
