package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"io"
	"text/tabwriter"
)

type folderNode struct {
	container
}

type folderCount struct {
	container
}

func (f *folderNode) LessThan(y interface{}) bool {
	return f.name < y.(*folderNode).name
}

func (f *folderNode) EqualTo(y interface{}) bool {
	return f.name == y.(*folderNode).name
}

func (f *folderCount) LessThan(y interface{}) bool {
	return f.count < y.(*folderCount).count
}

func (f *folderCount) EqualTo(y interface{}) bool {
	return f.count == y.(*folderCount).count
}

type moduleFolders struct {
	total    *totalInfo
	folders  *rbtree.RbTree
	topSize  *rbtree.RbTree
	topCount *rbtree.RbTree
}

type moduleFoldersNoOut struct {
	moduleFolders
}

func (m *moduleFolders) init() {
}

func (m *moduleFolders) postScan() {
	m.folders.WalkInorder(func(node *rbtree.Node) {
		fn := (*node.Key).(*folderNode)

		insertTo(m.topSize, &fn.container)

		fcn := folderCount{
			fn.container,
		}

		insertTo(m.topCount, &fcn)
	})

	m.total.CountFolders = m.folders.Root.Size
}

func (m *moduleFolders) folderHandler(fe *sys.FolderEntry) {
	var cmp rbtree.Comparable
	cmp = &folderNode{
		container{name: fe.Name, count: fe.Count, size: fe.Size},
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

	m.topSize.Descend(func(c *rbtree.Comparable) bool {

		folder := (*c).(*container)
		h := fmt.Sprintf("%d. %s", i, folder.name)

		i++

		count := folder.count
		sz := uint64(folder.size)

		outputTopStatLine(tw, count, m.total, sz, h)

		return true
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d folders by count:\n\n", top)
	_, _ = fmt.Fprintf(tw, format, "Folder", "Files", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "------", "-----", "------", "----", "------")

	i = 1

	m.topCount.Descend(func(c *rbtree.Comparable) bool {

		folder := (*c).(*folderCount)
		h := fmt.Sprintf("%d. %s", i, folder.name)

		i++

		count := folder.count
		sz := uint64(folder.size)

		outputTopStatLine(tw, count, m.total, sz, h)

		return true
	})

	_ = tw.Flush()
}
