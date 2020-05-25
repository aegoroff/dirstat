package cmd

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"io"
	"sync"
	"text/tabwriter"
	"time"
)

type module interface {
	handler() fileHandler
	output(tw *tabwriter.Writer, w io.Writer)
	postScan()
}

func executeModules(path string, fs afero.Fs, w io.Writer, fh folderHandler, modules []module) {
	var handlers []fileHandler
	for _, m := range modules {
		handlers = append(handlers, m.handler())
	}
	scan(path, fs, fh, handlers)

	for _, m := range modules {
		m.postScan()
	}

	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	for _, m := range modules {
		m.output(tw, w)
	}
}

func execute(path string, fs afero.Fs, w io.Writer, verbose bool, enabledRanges []int) {
	total := totalInfo{}

	var foldersMu sync.RWMutex
	folders := make(map[string]*container)

	// total module
	tm := moduleTotal{
		start: time.Now(),
		info:  &total,
	}

	// folders module
	fm := moduleFolders{
		&foldersMu,
		folders,
		rbtree.NewRbTree(),
		&total,
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

	// Aggregate module
	am := moduleAggregate{
		total:      &total,
		aggregator: make(map[string]countSizeAggregate),
	}

	// Total file stat module
	totfm := moduleTotalFile{
		total:     &total,
		aggregate: rangeAggregate,
	}

	// Modules order in the slice is important
	modules := []module{&totfm, &am, &tfm, &fm, &rm, &tm}

	foldersHandler := func(fsi *filesystemItem) {
		foldersMu.Lock()
		folders[fsi.dir] = &container{name: fsi.dir}
		total.CountFolders++
		foldersMu.Unlock()
	}

	executeModules(path, fs, w, foldersHandler, modules)
}
