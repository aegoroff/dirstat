package module

import (
	"errors"
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

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

type topper struct {
	headers []string
	p       out.Printer
	total   *totalInfo
}

func newTopper(p out.Printer, total *totalInfo, heads []string) *topper {
	return &topper{p: p, total: total, headers: heads}
}

func (t *topper) print(tree rbtree.RbTree, cast castFn) {
	tw := newTableWriter(t.p)

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

	it := rbtree.NewDescend(tree).Iterator()

	for it.Next() {
		fi, err := cast(it.Current())
		if err != nil {
			t.p.Cprint("<red>%v</>", err)
			return
		}

		count := fi.Count()
		sz := uint64(fi.Size())
		percentOfCount := t.total.countPercent(count)
		percentOfSize := t.total.sizePercent(sz)

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
