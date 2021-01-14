package module

import "github.com/aegoroff/godatastruct/rbtree"

// Context defines modules context
type Context struct {
	total *totalInfo
	top   int
	pd    decorator
}

// NewContext creates new module's context that needed to create new modules
func NewContext(top int, rr bool, root string) *Context {
	total := totalInfo{
		extensions: rbtree.New(),
	}

	var pd decorator

	if rr {
		pd = &removeRootDecorator{root: root}
	} else {
		pd = &noChangeDecorator{}
	}

	ctx := Context{
		total: &total,
		top:   top,
		pd:    pd,
	}
	return &ctx
}
