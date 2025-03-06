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

func (e *stringEnvironment) PrintFunc(w io.Writer, format string, a ...any) {
	_, _ = fmt.Fprint(w, e.SprintFunc(format, a...))
}

func (e *stringEnvironment) SprintFunc(format string, a ...any) string {
	s := fmt.Sprintf(format, a...)
	return e.re.ReplaceAllString(s, "$1")
}

func (e *stringEnvironment) Writer() io.WriteCloser {
	return e.w
}
