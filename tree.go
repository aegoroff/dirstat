package main

import "github.com/aegoroff/godatastruct/rbtree"

type fileNode struct {
	name string
	size int64
}

type folderNode struct {
	name  string
	size  int64
	count int64
}

type namedInts64 []*fileNode

func (x namedInts64) Len() int {
	return len(x)
}

func (x namedInts64) Less(i, j int) bool {
	return x[i].size < x[j].size
}

func (x namedInts64) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x *folderNode) LessThan(y interface{}) bool {
	return x.size < (y.(*folderNode)).size
}

func (x *folderNode) EqualTo(y interface{}) bool {
	return x.size == (y.(*folderNode)).size
}

func (x *fileNode) LessThan(y interface{}) bool {
	return x.size < (y.(*fileNode)).size
}

func (x *fileNode) EqualTo(y interface{}) bool {
	return x.size == (y.(*fileNode)).size
}

func newFolderTreeNode(f *folderNode) *rbtree.Comparable {
	var r rbtree.Comparable
	r = f
	return &r
}

func newFileTreeNode(f *fileNode) *rbtree.Comparable {
	var r rbtree.Comparable
	r = f
	return &r
}
