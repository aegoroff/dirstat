package module

import (
	"dirstat/module/internal/sys"
)

type fileHandler func(f *sys.FileEntry)

type fileFilterMiddleware struct {
	h fileHandler
}

func newFileFilterMiddleware(h fileHandler) *fileFilterMiddleware {
	return &fileFilterMiddleware{
		h: h,
	}
}

func (f *fileFilterMiddleware) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}
	f.h(evt.File)
}
