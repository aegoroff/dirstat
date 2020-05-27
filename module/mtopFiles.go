package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"io"
	"path/filepath"
	"text/tabwriter"
)

type moduleTopFiles struct {
	tree *rbtree.RbTree
}

type moduleTopFilesNoOut struct {
	moduleTopFiles
}

func (m *moduleTopFiles) init() {
}

func (m *moduleTopFiles) postScan() {

}

func (m *moduleTopFiles) folderHandler(_ *sys.FolderEntry) {

}

func (m *moduleTopFiles) fileHandler(f *sys.FileEntry) {
	fullPath := filepath.Join(f.Parent, f.Name)
	fileContainer := container{size: f.Size, name: fullPath, count: 1}
	fileContainer.insertTo(m.tree)
}

// Mute parent output
func (m *moduleTopFilesNoOut) output(_ *tabwriter.Writer, _ io.Writer) {

}

func (m *moduleTopFiles) output(tw *tabwriter.Writer, w io.Writer) {
	_, _ = fmt.Fprintf(w, "\nTOP %d files by size:\n\n", top)
	_, _ = fmt.Fprintf(tw, "%v\t%v\n", "File", "Size")
	_, _ = fmt.Fprintf(tw, "%v\t%v\n", "------", "----")

	i := 1

	m.tree.Descend(func(c *rbtree.Comparable) bool {
		file := (*c).(*container)
		h := fmt.Sprintf("%d. %s", i, file.name)

		i++

		sz := uint64(file.size)

		_, _ = fmt.Fprintf(tw, "%v\t%v\n", h, humanize.IBytes(sz))

		return true
	})

	_ = tw.Flush()
}
