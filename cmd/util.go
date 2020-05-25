package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"io"
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
	entry *fileInfo
	event fsEvent
}

type fileInfo struct {
	isDir bool
	name  string
	size  int64
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

func walkDirBreadthFirst(path string, fs afero.Fs, results chan<- *filesystemItem) {
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

			dirEvent := filesystemItem{
				dir:   d,
				event: fsEventDir,
			}
			results <- &dirEvent

			for _, entry := range entries {
				// Queue subdirs to walk in a queue
				if entry.isDir {
					subdir := filepath.Join(d, entry.name)

					// Push
					mu.Lock()
					queue = append(queue, subdir)
					mu.Unlock()
				} else {
					// Send to channel
					fileEvent := filesystemItem{
						dir:   d,
						entry: entry,
						event: fsEventFile,
					}
					results <- &fileEvent
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

var sema = make(chan struct{}, 32)

func dirents(path string, fs afero.Fs) []*fileInfo {
	sema <- struct{}{}
	defer func() { <-sema }()
	f, err := fs.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(-1)
	if err != nil {
		return nil
	}

	var result = []*fileInfo{}
	for _, e := range entries {
		fi := fileInfo{name: e.Name(), size: e.Size(), isDir: e.IsDir()}
		result = append(result, &fi)
	}

	return result
}
