package sys

import (
	"github.com/spf13/afero"
	"path/filepath"
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
	defer bf.closeRestrict()

	bf.push(path)

	for bf.len() > 0 {
		currentDir := bf.pop()

		bf.addOne()
		go bf.walk(currentDir, results)

		if bf.len() == 0 {
			// Waiting pending goroutines
			bf.wait()
		}
	}
}
