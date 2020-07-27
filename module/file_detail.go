package module

import (
	"dirstat/module/internal/sys"
	"sort"
)

type detailFileWorker struct {
	voidFinalize
	*fileFilter
	distribution     map[Range]files
	enabledRanges    []int
	enabledRangesMap map[int]bool
	fileRanges       ranges
	pd               *pathDecorator
}

type detailFileRenderer struct {
	*detailFileWorker
}

func newDetailFileWorker(rs ranges, enabledRanges []int, pd *pathDecorator) *detailFileWorker {
	w := detailFileWorker{
		enabledRanges: enabledRanges,
		distribution:  make(map[Range]files, len(rs)),
		fileRanges:    rs,
		pd:            pd,
	}

	w.fileFilter = newFileFilter(w.onFile)

	return &w
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

func (m *detailFileWorker) onFile(f *sys.FileEntry) {
	// Calculate files range statistic
	for i, r := range m.fileRanges {
		// Store each file info within range only if this this file size detail option set
		if !r.Contains(f.Size) || !m.enabledRangesMap[i+1] {
			continue
		}

		nodes, ok := m.distribution[r]
		if !ok {
			m.distribution[r] = make(files, 0)
		}
		fileContainer := file{size: f.Size, path: f.Path, pd: m.pd}
		m.distribution[r] = append(nodes, &fileContainer)
	}
}

// Renderer method

func (m *detailFileRenderer) print(p printer) {
	if len(m.enabledRanges) > 0 {
		heads := m.fileRanges.heads(true)
		p.cprint("\n<gray>Detailed files stat:</>\n")
		for i, r := range m.fileRanges {
			if len(m.distribution[r]) == 0 {
				continue
			}

			p.cprint("<gray>%s</>\n", heads[i])

			files := m.distribution[r]
			sort.Sort(sort.Reverse(files))

			for _, f := range files {
				size := human(f.size)
				p.cprint("   %s - <yellow>%s</>\n", f, size)
			}
		}
	}
}
