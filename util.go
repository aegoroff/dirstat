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

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage(w io.Writer) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Fprintf(w, "\nAlloc = %s", humanize.IBytes(m.Alloc))
	fmt.Fprintf(w,"\tTotalAlloc = %s", humanize.IBytes(m.TotalAlloc))
	fmt.Fprintf(w,"\tSys = %s", humanize.IBytes(m.Sys))
	fmt.Fprintf(w,"\tNumGC = %v\n", m.NumGC)
}

func walkDirBreadthFirst(path string, fs afero.Fs, action func(parent string, entry os.FileInfo)) {
	queue := make([]string, 0)

	queue = append(queue, path)

	for len(queue) > 0 {
		curr := queue[0]

		for _, entry := range dirents(curr, fs) {
			action(curr, entry)
			if entry.IsDir() {
				queue = append(queue, filepath.Join(curr, entry.Name()))
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
