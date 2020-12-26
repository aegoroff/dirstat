package module

import "dirstat/module/internal/sys"

type onlyFilesWorker struct {
	h sys.Handler
}

type onlyFoldersWorker struct {
	h sys.Handler
}

func newOnlyFilesWorker(h sys.Handler) sys.Handler {
	return &onlyFilesWorker{
		h: h,
	}
}

func newOnlyFoldersWorker(h sys.Handler) sys.Handler {
	return &onlyFoldersWorker{
		h: h,
	}
}

func (f *onlyFilesWorker) Handle(evt *sys.ScanEvent) {
	if evt.File != nil {
		f.h.Handle(evt)
	}
}

func (f *onlyFoldersWorker) Handle(evt *sys.ScanEvent) {
	if evt.Folder != nil {
		f.h.Handle(evt)
	}
}
