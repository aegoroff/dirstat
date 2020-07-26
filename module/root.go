package module

type rootMiddleware struct {
	removeRoot bool
	root       string
}

func (rm *rootMiddleware) decorate(r string) string {
	if rm.removeRoot {
		return r[len(rm.root):]
	}
	return r
}
