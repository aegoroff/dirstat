package module

import (
	"errors"
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// folder represents file system container that described by path
// and has size and the number of elements in it (count field).
type folder struct {
	path  string
	size  int64
	count int64
	pd    decorator
}

// Count sortable folder
type folderC struct {
	folder
}

// Size sortable folder
type folderS struct {
	folder
}

// Path sortable folder methods

func (f *folder) String() string { return f.Path() }
func (f *folder) Size() int64    { return f.size }
func (f *folder) Count() int64   { return f.count }
func (f *folder) Path() string   { return f.pd.decorate(f.path) }

// Count sortable folder methods

func (fc *folderC) Less(y rbtree.Comparable) bool  { return fc.count < y.(*folderC).count }
func (fc *folderC) Equal(y rbtree.Comparable) bool { return fc.count == y.(*folderC).count }

// Size sortable folder methods

func (fs *folderS) Less(y rbtree.Comparable) bool  { return fs.size < y.(*folderS).size }
func (fs *folderS) Equal(y rbtree.Comparable) bool { return fs.size == y.(*folderS).size }

func castSize(c rbtree.Comparable) (folderI, error) {
	f, ok := c.(*folderS)

	if !ok {
		return nil, errors.New("invalid casting: expected *folderS key type but it wasn`t")
	}
	return f, nil
}

func castCount(c rbtree.Comparable) (folderI, error) {
	f, ok := c.(*folderC)

	if !ok {
		return nil, errors.New("invalid casting: expected *folderC key type but it wasn`t")
	}
	return f, nil
}

type castFn func(c rbtree.Comparable) (folderI, error)

type top struct {
	tree    rbtree.RbTree
	cast    castFn
	headers []string
}

func newTop(tree rbtree.RbTree, cast castFn, heads []string) *top {
	return &top{tree: tree, cast: cast, headers: heads}
}

func (t *top) print(p out.Printer, total *totalInfo) {
	tw := newTableWriter(p)

	tw.addHeaders(t.headers)

	tw.configColumns([]table.ColumnConfig{
		{Number: 1, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 2, Align: text.AlignLeft, AlignHeader: text.AlignLeft, WidthMax: 100},
		{Number: 3, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.countTransformer},
		{Number: 4, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
		{Number: 5, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.sizeTransformer},
		{Number: 6, Align: text.AlignLeft, AlignHeader: text.AlignLeft, Transformer: tw.percentTransformer},
	})

	i := 1

	it := rbtree.NewDescend(t.tree).Iterator()

	for it.Next() {
		fi, err := t.cast(it.Current())
		if err != nil {
			p.Cprint("<red>%v</>", err)
			return
		}

		count := fi.Count()
		sz := uint64(fi.Size())
		percentOfCount := total.countPercent(count)
		percentOfSize := total.sizePercent(sz)

		tw.addRow([]interface{}{
			i,
			fi.Path(),
			count,
			percentOfCount,
			sz,
			percentOfSize,
		})

		i++
	}

	tw.render()
}
