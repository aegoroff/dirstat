package module

import (
	"github.com/aegoroff/dirstat/scan"
	c9s "github.com/aegoroff/godatastruct/collections"
	"sort"
)

type detailsFile struct {
	distribution map[Range]files
	fileRanges   ranges
}

type detailFileHandler struct {
	*detailsFile
	enabled c9s.IntHashSet
	pd      decorator
}

type detailFileRenderer struct {
	*detailsFile
	*baseRenderer
}

func newDetailsFile(rs ranges) *detailsFile {
	w := detailsFile{
		fileRanges:   rs,
		distribution: make(map[Range]files, len(rs)),
	}

	return &w
}

func newDetailFileHandler(details *detailsFile, enabledRanges []int, pd decorator) scan.Handler {
	er := make(c9s.IntHashSet)
	er.AddRange(enabledRanges...)
	w := detailFileHandler{
		detailsFile: details,
		enabled:     er,
		pd:          pd,
	}

	return newOnlyFilesHandler(&w)
}

func newDetailFileRenderer(details *detailsFile, order int) renderer {
	return &detailFileRenderer{
		detailsFile:  details,
		baseRenderer: newBaseRenderer(order),
	}
}

// Worker method

func (m *detailFileHandler) Handle(evt *scan.Event) {
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
	p.Cprint("\n<gray>Detailed files stat:</>\n")
	for i, r := range m.fileRanges {
		if len(m.distribution[r]) == 0 {
			continue
		}

		p.Cprint("<gray>%s</>\n", heads[i])

		files := m.distribution[r]
		sort.Sort(sort.Reverse(files))

		for _, f := range files {
			size := human(f.size)
			p.Cprint("   %s - <yellow>%s</>\n", f, size)
		}
	}
}
