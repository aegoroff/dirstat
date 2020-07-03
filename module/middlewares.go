package module

import (
	"dirstat/module/internal/sys"
)

type fileFilterMiddleware struct {
	voidInit
	voidFinalize
	wrk worker
}

func newFileFilterMiddleware(wrk worker) worker {
	return &fileFilterMiddleware{
		wrk: wrk,
	}
}

func (f *fileFilterMiddleware) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}
	f.wrk.handler(evt)
}
