package module

import "dirstat/module/internal/sys"

type onlyFilesWorker struct {
	w worker
}

func newOnlyFilesWorker(w worker) worker {
	return &onlyFilesWorker{
		w: w,
	}
}

func (f *onlyFilesWorker) init() {
	f.w.init()
}

func (f *onlyFilesWorker) finalize() {
	f.w.finalize()
}

func (f *onlyFilesWorker) handler(evt *sys.ScanEvent) {
	if evt.File != nil {
		f.w.handler(evt)
	}
}
