package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"io"
	"sync"
	"text/tabwriter"
	"time"
)

type module interface {
	handler() sys.FileHandler
	output(tw *tabwriter.Writer, w io.Writer)
	postScan()
}

// Execute runs all modules over path specified
func Execute(path string, fs afero.Fs, w io.Writer, verbose bool, enabledRanges []int) {
	total := totalInfo{}

	var foldersMu sync.RWMutex
	folders := make(map[string]*container)

	// total module
	totalm := moduleTotal{
		start: time.Now(),
		total: &total,
	}

	// folders module
	fm := moduleFolders{
		&foldersMu,
		&total,
		folders,
		rbtree.NewRbTree(),
	}

	// Top files module
	tfm := moduleTopFiles{
		tree: rbtree.NewRbTree(),
	}

	// verbose module
	verboseRanges := make(map[int]bool)
	for _, x := range enabledRanges {
		verboseRanges[x] = true
	}

	rangeAggregate := make(map[Range]fileStat)
	rm := moduleRange{
		verbose:       verbose,
		enabledRanges: verboseRanges,
		aggregate:     rangeAggregate,
		distribution:  make(map[Range]containers),
	}

	// File extensions statistic module
	em := moduleExtensions{
		total:      &total,
		aggregator: make(map[string]countSizeAggregate),
	}

	// Total file stat module
	totfm := moduleTotalFile{
		total:     &total,
		aggregate: rangeAggregate,
	}

	// Modules order in the slice is important
	modules := []module{&totfm, &em, &tfm, &fm, &rm, &totalm}

	foldersHandler := func(fsi *sys.FilesystemItem) {
		foldersMu.Lock()
		folders[fsi.Dir] = &container{name: fsi.Dir}
		total.CountFolders++
		foldersMu.Unlock()
	}

	executeModules(path, fs, w, foldersHandler, modules)
}

func executeModules(path string, fs afero.Fs, w io.Writer, fh sys.FolderHandler, modules []module) {
	var handlers []sys.FileHandler
	for _, m := range modules {
		handlers = append(handlers, m.handler())
	}
	sys.Scan(path, fs, fh, handlers)

	for _, m := range modules {
		m.postScan()
	}

	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	for _, m := range modules {
		m.output(tw, w)
	}
}

func outputTopStatLine(tw *tabwriter.Writer, count int64, total *totalInfo, sz uint64, title string) {
	percentOfCount := countPercent(count, total)
	percentOfSize := sizePercent(sz, total)

	_, _ = fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
}

func countPercent(count int64, total *totalInfo) float64 {
	return (float64(count) / float64(total.FilesTotal.Count)) * 100
}

func sizePercent(size uint64, total *totalInfo) float64 {
	return (float64(size) / float64(total.FilesTotal.Size)) * 100
}
