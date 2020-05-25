package cmd

import (
	"github.com/spf13/afero"
	"path/filepath"
	"sync"
)

type fsEvent int
type fileHandler func(f *fileEntry)
type folderHandler func(f *filesystemItem)

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

func scan(path string, fs afero.Fs, fh folderHandler, handlers []fileHandler) {
	filesystemCh := make(chan *filesystemItem, 1024)
	go func() {
		walkDirBreadthFirst(path, fs, filesystemCh)
	}()

	filesChan := make(chan *fileEntry, 1024)

	// Reading filesystem events
	go func() {
		defer close(filesChan)
		for item := range filesystemCh {
			if item.event == fsEventDir {
				fh(item)
			} else {
				// Only files
				entry := item.entry
				filesChan <- &fileEntry{Size: entry.size, Parent: item.dir, Name: entry.name}
			}
		}
	}()

	// Read all files from channel
	for file := range filesChan {
		for _, h := range handlers {
			h(file)
		}
	}
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
