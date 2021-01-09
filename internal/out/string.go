package out

import (
	"fmt"
	"io"
	"regexp"
)

type stringEnvironment struct {
	w  io.WriteCloser
	re *regexp.Regexp
}

func (e *stringEnvironment) NewPrinter() (Printer, error) {
	return NewPrinter(e), nil
}

func newStringEnvironment(w io.WriteCloser) *stringEnvironment {
	return &stringEnvironment{
		w:  w,
		re: regexp.MustCompile(`<[a-zA-Z_=,;]+>(.+?)</>`),
	}
}

func (e *stringEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	_, _ = fmt.Fprintf(w, e.re.ReplaceAllString(s, "$1"))
}

func (e *stringEnvironment) Writer() io.WriteCloser {
	return e.w
}
