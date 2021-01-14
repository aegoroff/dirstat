package module

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"time"
)

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
	if total == 0 {
		return 0
	}
	return (value / total) * 100
}

func numPrefixDecorator(ix int, h string) string {
	return fmt.Sprintf("%2d. %s", ix+1, h)
}

func human(n int64) string {
	return humanize.IBytes(uint64(n))
}
