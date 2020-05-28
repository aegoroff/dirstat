package sys

import (
	"github.com/spf13/afero"
	"path/filepath"
	"sync"
)

// ScanEvent defines scanning event structure
// that can contain file or folder event information
type ScanEvent struct {
	// File set not nil in case of file event occurred
	File *FileEntry

	// Folder set not nil in case of folder event occurred
	Folder *FolderEntry
}

// FileEntry represent file description
type FileEntry struct {
	// File size in bytes
	Size int64

	// Full path
	Path string
}

// FolderEntry represent folder description
type FolderEntry struct {
	FileEntry

	// The number of files in a folder
	Count int64
}

// FileHandler defines function prototype that handles each file event received
type ScanHandler func(f *ScanEvent)

type filesystemItem struct {
	dir   string
	name  string
	event fsEvent
	count int64
	size  int64
}

type filesysEntry struct {
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
func Scan(path string, fs afero.Fs, handlers []ScanHandler) {
	filesystemCh := make(chan *filesystemItem, 1024)
	go func() {
		walkDirBreadthFirst(path, fs, filesystemCh)
	}()

	scanChan := make(chan *ScanEvent, 1024)

	// Reading filesystem events
	go func() {
		defer close(scanChan)
		for item := range filesystemCh {
			se := ScanEvent{}
			if item.event == fsEventDir {
				fe := FileEntry{
					Size: item.size,
					Path: item.dir,
				}
				se.Folder = &FolderEntry{
					FileEntry: fe,
					Count:     item.count,
				}
			} else {
				se.File = &FileEntry{
					Size: item.size,
					Path: filepath.Join(item.dir, item.name),
				}
			}
			scanChan <- &se
		}
	}()

	// Read all files from channel
	for file := range scanChan {
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

			// Folder stat
			var count int64
			var size int64

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
						name:  entry.name,
						event: fsEventFile,
						count: 1,
						size:  entry.size,
					}
					results <- &fileEvent

					// update folder stat
					count++
					size += entry.size
				}
			}

			dirEvent := filesystemItem{
				dir:   d,
				event: fsEventDir,
				count: count,
				size:  size,
			}
			results <- &dirEvent

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

var concurrencyRestrictor = make(chan struct{}, 32)

func dirents(path string, fs afero.Fs) []*filesysEntry {
	concurrencyRestrictor <- struct{}{}
	defer func() { <-concurrencyRestrictor }()
	f, err := fs.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	entries, err := f.Readdir(-1)
	if err != nil {
		return nil
	}

	var result = []*filesysEntry{}
	for _, e := range entries {
		fi := filesysEntry{name: e.Name(), size: e.Size(), isDir: e.IsDir()}
		result = append(result, &fi)
	}

	return result
}
