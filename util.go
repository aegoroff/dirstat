package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type fsEvent int

const (
	fsEventDir  fsEvent = 0
	fsEventFile fsEvent = 1
)

type filesystemItem struct {
	dir   string
	entry os.FileInfo
	event fsEvent
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

func walkDirBreadthFirst(path string, fs afero.Fs, results chan<- filesystemItem) {
	defer close(results)

	var wg sync.WaitGroup
	var mu sync.RWMutex
	queue := make([]string, 0)

	queue = append(queue, path)

	ql := len(queue)

	for ql > 0 {
		// Peek
		mu.RLock()
		currentDir := queue[0]
		mu.RUnlock()

		wg.Add(1)
		go func(d string) {
			defer wg.Done()

			entries := dirents(d, fs)

			if entries == nil {
				return
			}

			results <- filesystemItem{
				dir:   d,
				event: fsEventDir,
			}

			for _, entry := range entries {
				// Queue subdirs to walk in a queue
				if entry.IsDir() {
					subdir := filepath.Join(d, entry.Name())

					// Push
					mu.Lock()
					queue = append(queue, subdir)
					mu.Unlock()
				} else {
					// Send to channel
					item := filesystemItem{
						dir:   d,
						entry: entry,
						event: fsEventFile,
					}
					results <- item
				}
			}
		}(currentDir)

		// Pop
		mu.Lock()
		queue = queue[1:]
		ql = len(queue)
		mu.Unlock()

		if ql == 0 {
			// Waiting pending goroutines
			wg.Wait()

			mu.RLock()
			ql = len(queue)
			mu.RUnlock()
		}
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
