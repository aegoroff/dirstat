package scan

// concurrentScans sets the default number of concurrent directory scans
const concurrentScans = 32

// Event defines scanning event structure
// that can contain file or folder event information.
// Each event can only be file or folder event but not
// file and folder simultaneously. Why do not make unified
// event type with three fields where count field in case
// of file will always be one? The answer is - memory. Typically
// there are much more files then folders in the filesystem so
// we can decrease memory consumption not to have extra 4 bytes
// for each file event because we know they are always contain one
type Event struct {
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
	path  string
	event fsEvent
	count int64
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

	scanEvents := make(chan *Event, 1024)

	go convert(fsEvents, scanEvents)

	// Read all files from channel
	for file := range scanEvents {
		for _, h := range handlers {
			h.Handle(file)
		}
	}
}

func convert(from <-chan *filesystemItem, to chan<- *Event) {
	defer close(to)
	for item := range from {
		se := Event{}
		if item.event == fsEventDir {
			se.Folder = newFolderEntry(item)
		} else {
			se.File = newFileEntry(item)
		}
		to <- &se
	}
}

func newFileEntry(item *filesystemItem) *FileEntry {
	return &FileEntry{
		Size: item.size,
		Path: item.path,
	}
}

func newFolderEntry(item *filesystemItem) *FolderEntry {
	return &FolderEntry{
		FileEntry: FileEntry{
			Size: item.size,
			Path: item.path,
		},
		Count: item.count,
	}
}

func walkBreadthFirst(path string, fs Filesystem, results chan<- *filesystemItem) {
	defer close(results)

	bf := newWalker(fs, concurrentScans)
	defer bf.closeRestrict()

	bf.enqueue(path)

	for bf.len() > 0 {
		currentDir := bf.dequeue()

		bf.addOne()
		go bf.walk(currentDir, results)

		if bf.len() == 0 {
			// Waiting pending goroutines
			bf.wait()
		}
	}
}
