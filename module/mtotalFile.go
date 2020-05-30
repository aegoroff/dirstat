package module

// NewTotalFileModule creates new total file statistic module
func NewTotalFileModule(ctx *Context) Module {
	r := &totalFileRenderer{
		total:     ctx.total,
		aggregate: ctx.rangeAggregate,
	}

	m := module{
		[]worker{},
		[]renderer{r},
	}

	return &m
}

type totalFileRenderer struct {
	total     *totalInfo
	aggregate map[Range]fileStat
}

func (m *totalFileRenderer) print(p printer) {
	p.print("Total files stat:\n\n")

	const format = "%v\t%v\t%v\t%v\t%v\n"

	p.printtab(format, "File size", "Amount", "%", "Size", "%")
	p.printtab(format, "---------", "------", "------", "----", "------")

	heads := createRangesHeads()
	for i, r := range fileSizeRanges {
		count := m.aggregate[r].TotalFilesCount
		sz := m.aggregate[r].TotalFilesSize

		m.total.printCountAndSizeStatLine(p, count, sz, heads[i])
	}
	p.flush()
}
