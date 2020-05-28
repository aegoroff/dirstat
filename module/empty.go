package module

import (
	"dirstat/module/internal/sys"
	"io"
	"text/tabwriter"
)

type emptyWorker struct{}
type emptyRenderer struct{}

func (m *emptyWorker) init() {}

func (m *emptyWorker) postScan() {}

func (m *emptyWorker) folderHandler(*sys.FolderEntry) {}

func (m *emptyWorker) fileHandler(*sys.FileEntry) {}

func (m *emptyRenderer) output(*tabwriter.Writer, io.Writer) {}
