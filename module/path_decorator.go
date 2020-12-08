package module

type decorator interface {
	decorate(s string) string
}

type nonDestructiveDecorator struct{}

type removeRootDecorator struct {
	root string
}

func (p *removeRootDecorator) decorate(s string) string { return s[len(p.root):] }

func (*nonDestructiveDecorator) decorate(s string) string { return s }
