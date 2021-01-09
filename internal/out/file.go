package out

import (
	"github.com/spf13/afero"
	"io"
)

type fileEnvironment struct {
	path *string
	fs   afero.Fs
	// base is wrapped PrintEnvironment
	// it may be memory buffer, console, etc.
	base PrintEnvironment
}

// NewWriteFileEnvironment creates new file output environment
func NewWriteFileEnvironment(path *string, fs afero.Fs, base PrintEnvironment) PrintEnvironment {
	return &fileEnvironment{
		path: path,
		fs:   fs,
		base: base,
	}
}

func (e *fileEnvironment) NewPrinter() (Printer, error) {
	// No path defined use base
	if *e.path == "" {
		return e.base.NewPrinter()
	}

	f, err := e.fs.Create(*e.path)
	if err != nil {
		return nil, err
	}

	e.base = newStringEnvironment(f)
	return NewPrinter(e), nil
}

func (e *fileEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	e.base.PrintFunc(w, format, a...)
}

func (e *fileEnvironment) Writer() io.WriteCloser {
	return e.base.Writer()
}
