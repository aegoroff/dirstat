package module

import (
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

var fileSizeRanges = [...]Range{
	{Min: 0, Max: 100 * kbyte},
	{Min: 100 * kbyte, Max: mbyte},
	{Min: mbyte, Max: 10 * mbyte},
	{Min: 10 * mbyte, Max: 100 * mbyte},
	{Min: 100 * mbyte, Max: gbyte},
	{Min: gbyte, Max: 10 * gbyte},
	{Min: 10 * gbyte, Max: 100 * gbyte},
	{Min: 100 * gbyte, Max: tbyte},
	{Min: tbyte, Max: 10 * tbyte},
	{Min: 10 * tbyte, Max: pbyte},
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

	p.printtab("%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
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
