package module

import (
	"dirstat/module/internal/sys"
)

type fileFilterMiddleware struct {
	wrk worker
}

func newFileFilterMiddleware(wrk worker) worker {
	return &fileFilterMiddleware{
		wrk: wrk,
	}
}

func (*fileFilterMiddleware) init() {}

func (f *fileFilterMiddleware) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}
	f.wrk.handler(evt)
}

func (*fileFilterMiddleware) finalize() {}
