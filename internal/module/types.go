package module

import (
	"fmt"
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

func (t *totalInfo) countPercent(count int64) float64 {
	return percent(float64(count), float64(t.FilesTotal.Count))
}

func (t *totalInfo) sizePercent(size uint64) float64 {
	return percent(float64(size), float64(t.FilesTotal.Size))
}

func (t *totalInfo) countExtensions() {
	t.CountFileExts = len(t.extensions)
}

func percent(value float64, total float64) float64 {
	return (value / total) * 100
}

type headDecorator func(ix int, h string) string

func transparentDecorator(_ int, h string) string {
	return h
}

func numPrefixDecorator(ix int, h string) string {
	return fmt.Sprintf("%2d. %s", ix+1, h)
}

func (r ranges) heads(hd headDecorator) []string {
	heads := make([]string, len(r))
	for i, r := range r {
		h := fmt.Sprintf("Between %s and %s", human(r.Min), human(r.Max))
		heads[i] = hd(i, h)
	}
	return heads
}
