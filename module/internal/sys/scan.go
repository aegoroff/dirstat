package sys

import (
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"sync"
)

// Handler defines scanning handler interface that handles filesystem events
type Handler interface {
	// Handle handles filesystem event
	Handle(evt *ScanEvent)
}

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

type breadthFirst struct {
	fs    afero.Fs
	wg    *sync.WaitGroup
	mu    *sync.RWMutex
	queue *[]string
	cr    chan struct{}
}

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
func Scan(path string, fs afero.Fs, handlers ...Handler) {
	fsEvents := make(chan *filesystemItem, 1024)
	go walkDirBreadthFirst(path, fs, fsEvents)

	scanEvents := make(chan *ScanEvent, 1024)

	go readFileSystemEvents(fsEvents, scanEvents)

	// Read all files from channel
	for file := range scanEvents {
		for _, h := range handlers {
			h.Handle(file)
		}
	}
}

func readFileSystemEvents(in <-chan *filesystemItem, out chan<- *ScanEvent) {
	defer close(out)
	for item := range in {
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
		out <- &se
	}
}

func walkDirBreadthFirst(path string, fs afero.Fs, results chan<- *filesystemItem) {
	defer close(results)
	var concurrencyRestrict = make(chan struct{}, 32)
	defer close(concurrencyRestrict)

	var wg sync.WaitGroup
	var mu sync.RWMutex
	queue := make([]string, 0)

	queue = append(queue, path)

	ql := len(queue)

	bf := &breadthFirst{
		fs:    fs,
		wg:    &wg,
		mu:    &mu,
		queue: &queue,
		cr:    concurrencyRestrict,
	}

	for ql > 0 {
		// Peek
		mu.RLock()
		currentDir := queue[0]
		mu.RUnlock()

		wg.Add(1)
		go bf.walk(currentDir, results)

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

func (bf *breadthFirst) walk(d string, results chan<- *filesystemItem) {
	defer bf.wg.Done()

	entries := bf.dirents(d)

	// Folder stat
	var count int64
	var size int64

	for _, entry := range entries {
		// Queue subdirs to walk in a queue
		if entry.isDir {
			subdir := filepath.Join(d, entry.name)

			// Push
			bf.mu.Lock()
			*bf.queue = append(*bf.queue, subdir)
			bf.mu.Unlock()
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
}

func (bf *breadthFirst) dirents(path string) []*filesysEntry {
	bf.cr <- struct{}{}
	defer func() { <-bf.cr }()
	f, err := bf.fs.Open(path)
	if err != nil {
		return []*filesysEntry{}
	}
	defer Close(f)

	entries, err := f.Readdir(-1)
	if err != nil {
		return []*filesysEntry{}
	}

	result := make([]*filesysEntry, 0, len(entries))
	for _, e := range entries {
		// dont follow symlinks
		if e.Mode()&os.ModeSymlink == 0 {
			fi := filesysEntry{name: e.Name(), size: e.Size(), isDir: e.IsDir()}
			result = append(result, &fi)
		}
	}

	return result
}
