package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

type filesystemItem struct {
	dir   string
	entry os.FileInfo
}

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage(w io.Writer) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	_, _ = fmt.Fprintf(w, "\nAlloc = %s", humanize.IBytes(m.Alloc))
	_, _ = fmt.Fprintf(w, "\tTotalAlloc = %s", humanize.IBytes(m.TotalAlloc))
	_, _ = fmt.Fprintf(w, "\tSys = %s", humanize.IBytes(m.Sys))
	_, _ = fmt.Fprintf(w, "\tNumGC = %v\n", m.NumGC)
}

//func walkDirBreadthFirst(path string, fs afero.Fs, action func(parent string, entry os.FileInfo)) {
func walkDirBreadthFirst(path string, fs afero.Fs, ch chan<- filesystemItem) {
	defer close(ch)
	queue := make([]string, 0)

	queue = append(queue, path)

	for len(queue) > 0 {
		curr := queue[0]

		for _, entry := range dirents(curr, fs) {
			// Send to channel
			item := filesystemItem{
				dir:   curr,
				entry: entry,
			}
			ch <- item

			// Queue subdirs to walk in a queue
			if entry.IsDir() {
				subdir := filepath.Join(curr, entry.Name())
				queue = append(queue, subdir)
			}
		}

		queue = queue[1:]
	}
}

func dirents(path string, fs afero.Fs) []os.FileInfo {
	f, err := fs.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(-1)
	if err != nil {
		return nil
	}

	return entries
}
