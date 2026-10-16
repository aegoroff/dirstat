package module

import (
	"os"
	"path/filepath"
)

type decorator interface {
	decorate(s string) string
}

// noChangeDecorator does nothing i.e. keep original string unchanged
type noChangeDecorator struct{}

type removeRootDecorator struct {
	root string
}

// decorate outputs path relative to root with leading separator.
// Root path itself becomes empty string
func (p *removeRootDecorator) decorate(s string) string {
	rel, err := filepath.Rel(p.root, s)
	if err != nil {
		return s
	}
	if rel == "." {
		return ""
	}
	return string(os.PathSeparator) + rel
}

func (*noChangeDecorator) decorate(s string) string { return s }
