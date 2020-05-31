package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"sort"
)

// Range defined integer value range
type Range struct {
	// Min value
	Min int64

	// Max value
	Max int64
}

// Contains defines whether the number specified within range
func (r *Range) Contains(num int64) bool {
	return num >= r.Min && num <= r.Max
}

type rangeWorker struct {
	distribution     map[Range]files
	aggregate        map[Range]fileStat
	verbose          bool
	enabledRanges    []int
	enabledRangesMap map[int]bool
}

type rangeRenderer struct {
	work *rangeWorker
}

func newRangeWorker(ctx *Context, verbose bool, enabledRanges []int) *rangeWorker {
	return &rangeWorker{
		verbose:       verbose,
		enabledRanges: enabledRanges,
		aggregate:     ctx.rangeAggregate,
		distribution:  make(map[Range]files),
	}
}

func newRangeRenderer(work *rangeWorker) renderer {
	return &rangeRenderer{work}
}

// Worker methods

func (m *rangeWorker) init() {
	m.enabledRangesMap = make(map[int]bool)
	for _, x := range m.enabledRanges {
		m.enabledRangesMap[x] = true
	}
}

func (m *rangeWorker) finalize() {

}

func (m *rangeWorker) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}
	f := evt.File

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
			m.distribution[r] = make(files, 0)
		}
		fileContainer := file{size: f.Size, path: f.Path}
		m.distribution[r] = append(nodes, &fileContainer)
	}
}

// Renderer method

func (m *rangeRenderer) print(p printer) {
	if m.work.verbose && len(m.work.enabledRanges) > 0 {
		heads := createRangesHeads()
		p.print("\nDetailed files stat:\n")
		for i, r := range fileSizeRanges {
			if len(m.work.distribution[r]) == 0 {
				continue
			}

			p.print("%s\n", heads[i])

			files := m.work.distribution[r]
			sort.Sort(sort.Reverse(files))

			for _, f := range files {
				size := human(f.size)
				p.print("   %s - %s\n", f, size)
			}
		}
	}
}

func createRangesHeads() []string {
	var heads []string
	for i, r := range fileSizeRanges {
		h := fmt.Sprintf("%d. Between %s and %s", i+1, human(r.Min), human(r.Max))
		heads = append(heads, h)
	}
	return heads
}
