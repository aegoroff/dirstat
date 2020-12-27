package module

import "github.com/aegoroff/dirstat/scan"

type onlyFilesHandler struct {
	h scan.Handler
}

type onlyFoldersHandler struct {
	h scan.Handler
}

func newOnlyFilesHandler(h scan.Handler) scan.Handler {
	return &onlyFilesHandler{
		h: h,
	}
}

func newOnlyFoldersHandler(h scan.Handler) scan.Handler {
	return &onlyFoldersHandler{
		h: h,
	}
}

func (f *onlyFilesHandler) Handle(evt *scan.ScanEvent) {
	if evt.File != nil {
		f.h.Handle(evt)
	}
}

func (f *onlyFoldersHandler) Handle(evt *scan.ScanEvent) {
	if evt.Folder != nil {
		f.h.Handle(evt)
	}
}
