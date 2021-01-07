package scan

import "os"

type osFs struct{}

// NewOsFs creates new real os backed Filesystem instance
func NewOsFs() Filesystem {
	return &osFs{}
}

func (o *osFs) Open(name string) (File, error) {
	return os.Open(name)
}
