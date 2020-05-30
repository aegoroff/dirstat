package module

import (
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

const top = 10

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
