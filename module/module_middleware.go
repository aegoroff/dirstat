package module

import "github.com/aegoroff/dirstat/module/internal/sys"

type onlyFilesHandler struct {
	h sys.Handler
}

type onlyFoldersHandler struct {
	h sys.Handler
}

func newOnlyFilesHandler(h sys.Handler) sys.Handler {
	return &onlyFilesHandler{
		h: h,
	}
}

func newOnlyFoldersHandler(h sys.Handler) sys.Handler {
	return &onlyFoldersHandler{
		h: h,
	}
}

func (f *onlyFilesHandler) Handle(evt *sys.ScanEvent) {
	if evt.File != nil {
		f.h.Handle(evt)
	}
}

func (f *onlyFoldersHandler) Handle(evt *sys.ScanEvent) {
	if evt.Folder != nil {
		f.h.Handle(evt)
	}
}
