package out

import (
	"fmt"
	"io"
)

// PrintEnvironment represents concrete printing environment abstraction
type PrintEnvironment interface {
	Writable

	// PrintFunc represents printing function implementation
	PrintFunc(w io.Writer, format string, a ...any)

	// SprintFunc represents printing to string function implementation
	SprintFunc(format string, a ...any) string

	// NewPrinter creates new printer
	NewPrinter() (Printer, error)
}

// StringEnvironment defines in memory printing environment abstraction
type StringEnvironment interface {
	PrintEnvironment
	fmt.Stringer
}

// Printer represents printing abstraction with colorizing support
type Printer interface {
	Writable
	// Cprint prints data with colorizing support
	Cprint(format string, a ...any)
	// Sprintf writes formatted string into another and do cleanup if necessary
	Sprintf(format string, a ...any) string
	// Println prints new line
	Println()
}

// Writable represents io.Writer container
type Writable interface {
	// Writer gets underlying io.Writer
	Writer() io.WriteCloser
}
