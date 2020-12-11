package module

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"time"
)

const (
	_ int64 = 1 << (10 * iota)
	kbyte
	mbyte
	gbyte
	tbyte
	pbyte
)

// Range defined integer value range
type Range struct {
	// Min value
	Min int64

	// Max value
	Max int64
}

type ranges []Range

// Contains defines whether the number specified within range
func (r *Range) Contains(num int64) bool {
	return num >= r.Min && num <= r.Max
}

type fileStat struct {
	TotalFilesSize  uint64
	TotalFilesCount int64
}

type totalInfo struct {
	ReadingTime   time.Duration
	FilesTotal    countSizeAggregate
	CountFolders  int64
	CountFileExts int
	extensions    map[string]countSizeAggregate
}

type countSizeAggregate struct {
	Count int64
	Size  uint64
}

type fixedTree struct {
	tree rbtree.RbTree
	size int
}

func (t *totalInfo) countPercent(count int64) float64 {
	return percent(float64(count), float64(t.FilesTotal.Count))
}

func (t *totalInfo) sizePercent(size uint64) float64 {
	return percent(float64(size), float64(t.FilesTotal.Size))
}

func percent(value float64, total float64) float64 {
	return (value / total) * 100
}

func newFixedTree(sz int) *fixedTree {
	return &fixedTree{
		tree: rbtree.NewRbTree(),
		size: sz,
	}
}

// insert inserts node into tree which size is limited
// Only <size> max nodes will be in the tree
func (t *fixedTree) insert(c rbtree.Comparable) {
	min := t.tree.Minimum()
	if t.tree.Len() < int64(t.size) || min.Key().LessThan(c) {
		if t.tree.Len() == int64(t.size) {
			t.tree.DeleteNode(min.Key())
		}

		t.tree.Insert(c)
	}
}

// descend walks tree in descending order
func (t *fixedTree) descend(callback rbtree.NodeEvaluator) {
	rbtree.NewDescend(t.tree).Foreach(callback)
}

type headDecorator func(ix int, h string) string

func transparentDecorator(_ int, h string) string {
	return h
}

func numPrefixDecorator(ix int, h string) string {
	return fmt.Sprintf("%2d. %s", ix+1, h)
}

func (r ranges) heads(hd headDecorator) []string {
	var heads []string
	for i, r := range r {
		h := fmt.Sprintf("Between %s and %s", human(r.Min), human(r.Max))
		heads = append(heads, hd(i, h))
	}
	return heads
}
