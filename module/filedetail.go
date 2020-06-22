package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"sort"
)

type detailFileWorker struct {
	distribution     map[Range]files
	enabledRanges    []int
	enabledRangesMap map[int]bool
}

type detailFileRenderer struct {
	work *detailFileWorker
}

func newDetailFileWorker(enabledRanges []int) *detailFileWorker {
	return &detailFileWorker{
		enabledRanges: enabledRanges,
		distribution:  make(map[Range]files),
	}
}

func newDetailFileRenderer(work *detailFileWorker) renderer {
	return &detailFileRenderer{work}
}

// Worker methods

func (m *detailFileWorker) init() {
	m.enabledRangesMap = make(map[int]bool)
	for _, x := range m.enabledRanges {
		m.enabledRangesMap[x] = true
	}
}

func (*detailFileWorker) finalize() {}

func (m *detailFileWorker) handler(evt *sys.ScanEvent) {
	if evt.File == nil {
		return
	}
	f := evt.File

	// Calculate files range statistic
	for i, r := range fileSizeRanges {
		// Store each file info within range only i verbose option set
		if !r.Contains(f.Size) || !m.enabledRangesMap[i+1] {
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

func (m *detailFileRenderer) print(p printer) {
	if len(m.work.enabledRanges) > 0 {
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
		h := fmt.Sprintf("%2d. Between %s and %s", i+1, human(r.Min), human(r.Max))
		heads = append(heads, h)
	}
	return heads
}
