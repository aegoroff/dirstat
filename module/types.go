package module

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
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

type ranges [10]Range

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
}

type countSizeAggregate struct {
	Count int64
	Size  uint64
}

func (t *totalInfo) countPercent(count int64) float64 {
	return (float64(count) / float64(t.FilesTotal.Count)) * 100
}

func (t *totalInfo) sizePercent(size uint64) float64 {
	return (float64(size) / float64(t.FilesTotal.Size)) * 100
}

func (t *totalInfo) printCountAndSizeStatLine(p printer, count int64, sz uint64, title string) {
	percentOfCount := t.countPercent(count)
	percentOfSize := t.sizePercent(sz)

	p.print("%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
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

func (r ranges) heads() []string {
	var heads []string
	for i, r := range r {
		h := fmt.Sprintf("%2d. Between %s and %s", i+1, human(r.Min), human(r.Max))
		heads = append(heads, h)
	}
	return heads
}
