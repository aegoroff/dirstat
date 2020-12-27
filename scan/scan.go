package scan

import (
	"io"
	"os"
	"path/filepath"
)

// concurrentScans sets the default number of concurrent directory scans
const concurrentScans = 32

// Handler defines scanning handler interface that handles filesystem events
type Handler interface {
	// Handle handles filesystem event
	Handle(evt *ScanEvent)
}

// Filesystem represents filesystem abstraction
type Filesystem interface {
	// Open opens file for reading
	Open(name string) (File, error)
}

// File represents an open file descriptor.
type File interface {
	io.Closer

	// Readdir reads the contents of the directory associated with file and
	// returns a slice of up to n FileInfo values, as would be returned
	// by Lstat, in directory order. Subsequent calls on the same file will yield
	// further FileInfos.
	// If n > 0, Readdir returns at most n FileInfo structures. In this case, if
	// Readdir returns an empty slice, it will return a non-nil error
	// explaining why. At the end of a directory, the error is io.EOF.
	//
	// If n <= 0, Readdir returns all the FileInfo from the directory in
	// a single slice. In this case, if Readdir succeeds (reads all
	// the way to the end of the directory), it returns the slice and a
	// nil error. If it encounters an error before the end of the
	// directory, Readdir returns the FileInfo read until that point
	// and a non-nil error.
	Readdir(count int) ([]os.FileInfo, error)
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
func Scan(path string, fs Filesystem, handlers ...Handler) {
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

func walkBreadthFirst(path string, fs Filesystem, results chan<- *filesystemItem) {
	defer close(results)

	bf := newWalker(fs, concurrentScans)
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
