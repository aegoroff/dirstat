package out

import (
	"io"
	"os"

	"github.com/gookit/color"
)

type consoleEnvironment struct {
	stdout *stdout
}

type stdout struct {
	w io.Writer
}

func (s *stdout) Write(p []byte) (n int, err error) {
	return s.w.Write(p)
}

func (*stdout) Close() error {
	return nil
}

// NewConsoleEnvironment creates new PrintEnvironment that outputs to console
func NewConsoleEnvironment() PrintEnvironment {
	return &consoleEnvironment{stdout: &stdout{w: os.Stdout}}
}

func (e *consoleEnvironment) NewPrinter() (Printer, error) {
	return NewPrinter(e), nil
}

func (*consoleEnvironment) PrintFunc(w io.Writer, format string, a ...any) {
	color.Fprintf(w, format, a...)
}

func (*consoleEnvironment) SprintFunc(format string, a ...any) string {
	return color.Sprintf(format, a...)
}

func (e *consoleEnvironment) Writer() io.WriteCloser {
	return e.stdout
}
