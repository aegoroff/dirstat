package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"io"
	"text/tabwriter"
)

type folderNode struct {
	path  string
	value *container
}

func (f *folderNode) LessThan(y interface{}) bool {
	return f.path < y.(*folderNode).path
}

func (f *folderNode) EqualTo(y interface{}) bool {
	return f.path == y.(*folderNode).path
}

type moduleFolders struct {
	total   *totalInfo
	folders *rbtree.RbTree
	top     *rbtree.RbTree
}

type moduleFoldersNoOut struct {
	moduleFolders
}

func (m *moduleFolders) init() {
}

func (m *moduleFolders) postScan() {
	m.folders.WalkInorder(func(node *rbtree.Node) {
		n := (*node.Key).(*folderNode)
		cont := n.value
		cont.insertTo(m.top)
	})

	m.total.CountFolders = m.folders.Root.Size
}

func (m *moduleFolders) folderHandler(fe *sys.FolderEntry) {
	var cmp rbtree.Comparable
	cmp = &folderNode{
		path:  fe.Name,
		value: &container{name: fe.Name, count: fe.Count, size: fe.Size},
	}
	m.folders.Insert(rbtree.NewNode(&cmp))
}

func (m *moduleFolders) fileHandler(_ *sys.FileEntry) {

}

// Mute parent output
func (m *moduleFoldersNoOut) output(_ *tabwriter.Writer, _ io.Writer) {

}

func (m *moduleFolders) output(tw *tabwriter.Writer, w io.Writer) {
	const format = "%v\t%v\t%v\t%v\t%v\n"

	_, _ = fmt.Fprintf(w, "\nTOP %d folders by size:\n\n", top)
	_, _ = fmt.Fprintf(tw, format, "Folder", "Files", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "------", "-----", "------", "----", "------")

	i := 1

	m.top.Descend(func(c *rbtree.Comparable) bool {

		folder := (*c).(*container)
		h := fmt.Sprintf("%d. %s", i, folder.name)

		i++

		count := folder.count
		sz := uint64(folder.size)

		outputTopStatLine(tw, count, m.total, sz, h)

		return true
	})

	_ = tw.Flush()
}
