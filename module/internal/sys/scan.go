package sys

import (
	"github.com/spf13/afero"
	"path/filepath"
	"sync"
)

type FileEntry struct {
	Size   int64
	Parent string
	Name   string
}
type FileHandler func(f *FileEntry)
type FolderHandler func(f *FilesystemItem)

// FilesystemItem defines filesystem abstraction (file or folder)
type FilesystemItem struct {
	// Item's dir
	Dir string

	// Items filesystem info like size, name, etc.
	Entry *FileInfo
	event fsEvent
}

// FileInfo defines filesystem item data like size, name, etc.
type FileInfo struct {
	isDir bool
	name  string
	size  int64
}

type fsEvent int

const (
	fsEventDir  fsEvent = 0
	fsEventFile fsEvent = 1
)

// Scan do specified path scanning and executes folder handler on each folder
// and all file handlers on each file
func Scan(path string, fs afero.Fs, fh FolderHandler, handlers []FileHandler) {
	filesystemCh := make(chan *FilesystemItem, 1024)
	go func() {
		walkDirBreadthFirst(path, fs, filesystemCh)
	}()

	filesChan := make(chan *FileEntry, 1024)

	// Reading filesystem events
	go func() {
		defer close(filesChan)
		for item := range filesystemCh {
			if item.event == fsEventDir {
				fh(item)
			} else {
				// Only files
				entry := item.Entry
				filesChan <- &FileEntry{Size: entry.size, Parent: item.Dir, Name: entry.name}
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

func walkDirBreadthFirst(path string, fs afero.Fs, results chan<- *FilesystemItem) {
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

			dirEvent := FilesystemItem{
				Dir:   d,
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
					fileEvent := FilesystemItem{
						Dir:   d,
						Entry: entry,
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

func dirents(path string, fs afero.Fs) []*FileInfo {
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

	var result = []*FileInfo{}
	for _, e := range entries {
		fi := FileInfo{name: e.Name(), size: e.Size(), isDir: e.IsDir()}
		result = append(result, &fi)
	}

	return result
}
