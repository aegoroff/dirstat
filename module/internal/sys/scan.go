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

type walker struct {
	fs         afero.Fs
	wg         sync.WaitGroup
	mu         sync.RWMutex
	queue      []string
	restrictor chan struct{}
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
	go walkBreadthFirst(path, fs, fsEvents)

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
			se.Folder = newFolderEntry(item)
		} else {
			se.File = newFileEntry(item)
		}
		out <- &se
	}
}

func newFileEntry(item *filesystemItem) *FileEntry {
	return &FileEntry{
		Size: item.size,
		Path: filepath.Join(item.dir, item.name),
	}
}

func newFolderEntry(item *filesystemItem) *FolderEntry {
	return &FolderEntry{
		FileEntry: FileEntry{
			Size: item.size,
			Path: item.dir,
		},
		Count: item.count,
	}
}

func walkBreadthFirst(path string, fs afero.Fs, results chan<- *filesystemItem) {
	defer close(results)

	bf := newWalker(fs, 32)
	defer close(bf.restrictor)

	bf.push(path)

	ql := len(bf.queue)

	for ql > 0 {
		currentDir := bf.peek()

		bf.wg.Add(1)
		go bf.walk(currentDir, results)

		ql = bf.pop()

		if ql == 0 {
			// Waiting pending goroutines
			bf.wait()

			ql = bf.len()
		}
	}
}

func newWalker(fs afero.Fs, parallel int) *walker {
	bf := &walker{
		fs:         fs,
		queue:      make([]string, 0),
		restrictor: make(chan struct{}, parallel),
	}
	return bf
}

func (bf *walker) peek() string {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	return bf.queue[0]
}

func (bf *walker) pop() int {
	bf.mu.Lock()
	bf.mu.Unlock()
	bf.queue = bf.queue[1:]
	return len(bf.queue)
}

func (bf *walker) push(s string) {
	bf.mu.Lock()
	bf.queue = append(bf.queue, s)
	bf.mu.Unlock()
}

func (bf *walker) len() int {
	bf.mu.RLock()
	defer bf.mu.RUnlock()
	return len(bf.queue)
}

func (bf *walker) wait() {
	bf.wg.Wait()
}

func (bf *walker) walk(d string, results chan<- *filesystemItem) {
	defer bf.wg.Done()

	entries := bf.readdir(d)

	// Folder stat
	var count int64
	var size int64

	for _, entry := range entries {
		// Queue subdirs to walk in a queue
		if entry.isDir {
			bf.push(filepath.Join(d, entry.name))
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

func (bf *walker) readdir(path string) []*filesysEntry {
	bf.restrictor <- struct{}{}
	defer func() { <-bf.restrictor }()
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
