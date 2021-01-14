package module

type decorator interface {
	decorate(s string) string
}

// noChangeDecorator does nothing i.e. keep original string unchanged
type noChangeDecorator struct{}

type removeRootDecorator struct {
	root string
}

func (p *removeRootDecorator) decorate(s string) string { return s[len(p.root):] }

func (*noChangeDecorator) decorate(s string) string { return s }
