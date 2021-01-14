package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"sort"
)

type detailsFile struct {
	enabled []*Range
}

type detailFileHandler struct {
	*detailsFile
}

type detailFileRenderer struct {
	*detailsFile
	*baseRenderer
	pd decorator
}

func newDetailsFile(rs rbtree.RbTree, enabledRanges []int) *detailsFile {
	w := &detailsFile{
		enabled: make([]*Range, 0, len(enabledRanges)),
	}

	for _, er := range enabledRanges {
		n, ok := rs.OrderStatisticSelect(int64(er))
		if ok {
			w.enabled = append(w.enabled, n.Key().(*Range))
		}
	}

	return w
}

func newDetailFileHandler(details *detailsFile) scan.Handler {
	w := &detailFileHandler{details}
	return newOnlyFilesHandler(w)
}

func newDetailFileRenderer(details *detailsFile, pd decorator, order int) renderer {
	return &detailFileRenderer{
		detailsFile:  details,
		baseRenderer: newBaseRenderer(order),
		pd:           pd,
	}
}

// Worker method

func (m *detailFileHandler) Handle(evt *scan.Event) {
	f := evt.File

	// Calculate files range statistic
	for _, r := range m.enabled {
		// Store each file info within range only if this this file size detail option set
		if !r.Contains(f.Size) {
			continue
		}

		f := &file{size: f.Size, path: f.Path}
		r.files = append(r.files, f)
	}
}

// Renderer method

func (m *detailFileRenderer) render(p out.Printer) {
	p.Cprint("\n<gray>Detailed files stat:</>\n")
	for i, r := range m.enabled {
		if len(r.files) == 0 {
			continue
		}

		p.Cprint("<gray>%s</>\n", numPrefixDecorator(i, r.String()))

		sort.Sort(sort.Reverse(r.files))

		for _, f := range r.files {
			size := human(f.size)
			p.Cprint("   %s - <yellow>%s</>\n", f, size)
		}
	}
}
