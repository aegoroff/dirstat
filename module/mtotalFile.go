package module

// NewTotalFileModule creates new total file statistic module
func NewTotalFileModule(ctx *Context) Module {
	m := moduleTotalFile{
		total:     ctx.total,
		aggregate: ctx.rangeAggregate,
	}
	return &m
}

type moduleTotalFile struct {
	emptyWorker
	total     *totalInfo
	aggregate map[Range]fileStat
}

func (m *moduleTotalFile) output(p printer) {
	p.print("Total files stat:\n\n")

	const format = "%v\t%v\t%v\t%v\t%v\n"

	p.printtab(format, "File size", "Amount", "%", "Size", "%")
	p.printtab(format, "---------", "------", "------", "----", "------")

	heads := createRangesHeads()
	for i, r := range fileSizeRanges {
		count := m.aggregate[r].TotalFilesCount
		sz := m.aggregate[r].TotalFilesSize

		m.total.printTopStatLine(p, count, sz, heads[i])
	}
	p.flush()
}
