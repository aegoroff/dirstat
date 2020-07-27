package module

type pathDecorator struct {
	removeRoot bool
	root       string
}

func (p *pathDecorator) decorate(full string) string {
	if p.removeRoot {
		return full[len(p.root):]
	}
	return full
}
