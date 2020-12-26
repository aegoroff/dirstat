package module

import (
	"dirstat/module/internal/sys"
	c9s "github.com/aegoroff/godatastruct/collections"
	"sort"
)

type detailFileHandler struct {
	distribution map[Range]files
	enabled      c9s.IntHashSet
	fileRanges   ranges
	pd           decorator
}

type detailFileRenderer struct {
	*detailFileHandler
	*baseRenderer
}

func newDetailFileHandler(rs ranges, enabledRanges []int, pd decorator) *detailFileHandler {
	er := make(c9s.IntHashSet)
	er.AddRange(enabledRanges...)
	w := detailFileHandler{
		enabled:      er,
		distribution: make(map[Range]files, len(rs)),
		fileRanges:   rs,
		pd:           pd,
	}

	return &w
}

func newDetailFileRenderer(work *detailFileHandler, order int) renderer {
	return &detailFileRenderer{
		detailFileHandler: work,
		baseRenderer:      newBaseRenderer(order),
	}
}

// Worker method

func (m *detailFileHandler) Handle(evt *sys.ScanEvent) {
	f := evt.File
	// Calculate files range statistic
	for i, r := range m.fileRanges {
		// Store each file info within range only if this this file size detail option set
		if !r.Contains(f.Size) || !m.enabled.Contains(i+1) {
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
	heads := m.fileRanges.heads(numPrefixDecorator)
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
