package module

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"time"
)

type totalInfo struct {
	ReadingTime     time.Duration
	FilesTotal      int64
	FilesSize       int64
	FoldersTotal    int64
	ExtensionsTotal int64
	extensions      rbtree.RbTree
}

func (t *totalInfo) countPercent(count int64) float64 {
	return percent(float64(count), float64(t.FilesTotal))
}

func (t *totalInfo) sizePercent(size int64) float64 {
	return percent(float64(size), float64(t.FilesSize))
}

func percent(value float64, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (value / total) * 100
}

func humanSize(n int64) string {
	return humanize.IBytes(uint64(n))
}
