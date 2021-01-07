package scan

import (
	"io"
	"os"
)

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
