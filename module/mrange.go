package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"sort"
	"text/tabwriter"
)

// Range defined integer value range
type Range struct {
	// Min value
	Min int64

	// Max value
	Max int64
}

// Contains defines whether the number specified within range
func (r Range) Contains(num int64) bool {
	return num >= r.Min && num <= r.Max
}

// NewRangeModule creates new file statistic by file size range module
func NewRangeModule(ctx *Context, verbose bool, enabledRanges []int) Module {
	work := newRangeWorker(ctx, verbose, enabledRanges)
	rend := rangeRenderer{work}
	m := moduleRange{
		work,
		rend,
	}
	return &m
}

// NewRangeHiddenModule creates new file statistic by file size range module
// that has disabled output
func NewRangeHiddenModule(ctx *Context) Module {
	work := newRangeWorker(ctx, false, []int{})

	m := moduleRangeNoOut{
		work,
		emptyRenderer{},
	}
	return &m
}

type rangeWorker struct {
	distribution     map[Range]containers
	aggregate        map[Range]fileStat
	verbose          bool
	enabledRanges    []int
	enabledRangesMap map[int]bool
}

type rangeRenderer struct {
	rangeWorker
}

type moduleRange struct {
	rangeWorker
	rangeRenderer
}

type moduleRangeNoOut struct {
	rangeWorker
	emptyRenderer
}

func newRangeWorker(ctx *Context, verbose bool, enabledRanges []int) rangeWorker {
	return rangeWorker{
		verbose:       verbose,
		enabledRanges: enabledRanges,
		aggregate:     ctx.rangeAggregate,
		distribution:  make(map[Range]containers),
	}
}

func (m *rangeWorker) init() {
	m.enabledRangesMap = make(map[int]bool)
	for _, x := range m.enabledRanges {
		m.enabledRangesMap[x] = true
	}
}

func (m *rangeWorker) postScan() {

}

func (m *rangeWorker) folderHandler(_ *sys.FolderEntry) {

}

func (m *rangeWorker) fileHandler(f *sys.FileEntry) {
	unsignedSize := uint64(f.Size)

	// Calculate files range statistic
	for i, r := range fileSizeRanges {
		if !r.Contains(f.Size) {
			continue
		}

		s := m.aggregate[r]
		s.TotalFilesCount++
		s.TotalFilesSize += unsignedSize
		m.aggregate[r] = s

		// Store each file info within range only i verbose option set
		if !m.verbose || !m.enabledRangesMap[i+1] {
			continue
		}

		nodes, ok := m.distribution[r]
		if !ok {
			m.distribution[r] = make(containers, 0)
		}
		fileContainer := container{size: f.Size, name: f.Path, count: 1}
		m.distribution[r] = append(nodes, &fileContainer)
	}
}

func (m *rangeRenderer) output(_ *tabwriter.Writer, w io.Writer) {
	if m.verbose && len(m.enabledRanges) > 0 {
		heads := createRangesHeads()
		_, _ = fmt.Fprintf(w, "\nDetailed files stat:\n")
		for i, r := range fileSizeRanges {
			if len(m.distribution[r]) == 0 {
				continue
			}

			_, _ = fmt.Fprintf(w, "%s\n", heads[i])

			items := m.distribution[r]
			sort.Sort(sort.Reverse(items))

			for _, item := range items {
				size := humanize.IBytes(uint64(item.size))
				_, _ = fmt.Fprintf(w, "   %s - %s\n", item.name, size)
			}
		}
	}
}

func createRangesHeads() []string {
	var heads []string
	for i, r := range fileSizeRanges {
		h := fmt.Sprintf("%d. Between %s and %s", i+1, humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)))
		heads = append(heads, h)
	}
	return heads
}
