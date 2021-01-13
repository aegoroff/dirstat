package out

import (
	"github.com/gookit/color"
	"io"
	"os"
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

func (*consoleEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	color.Fprintf(w, format, a...)
}

func (e *consoleEnvironment) SprintFunc(format string, a ...interface{}) string {
	return color.Sprintf(format, a...)
}

func (e *consoleEnvironment) Writer() io.WriteCloser {
	return e.stdout
}
