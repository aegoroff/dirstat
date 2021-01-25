package scan

import "os"

type osFs struct{}

// NewOsFs creates new real os backed Filesystem instance
func NewOsFs() Filesystem {
	return &osFs{}
}

// Open opens the named file for reading.
// If successful, methods on the returned file can be used for reading;
// the associated file descriptor has mode os.O_RDONLY.
// If there is an error, it will be of type *os.PathError.
func (*osFs) Open(name string) (File, error) {
	return os.Open(name)
}
