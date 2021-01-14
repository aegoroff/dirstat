package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"sort"
)

type detailsFile struct {
	enabled map[int]*Range
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
		enabled: make(map[int]*Range, len(enabledRanges)),
	}

	for _, er := range enabledRanges {
		n, ok := rs.OrderStatisticSelect(int64(er))
		if ok {
			w.enabled[er] = n.Key().(*Range)
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

	keys := m.enabledKeys()

	for _, k := range keys {
		r := m.enabled[k]
		if len(r.files) == 0 {
			continue
		}

		p.Cprint("<gray>%2d. %s</>\n", k, r.String())

		sort.Sort(sort.Reverse(r.files))

		for _, f := range r.files {
			size := human(f.size)
			p.Cprint("   %s - <yellow>%s</>\n", f, size)
		}
	}
}

func (m *detailFileRenderer) enabledKeys() sort.IntSlice {
	keys := make(sort.IntSlice, len(m.enabled))

	i := 0
	for k := range m.enabled {
		keys[i] = k
		i++
	}
	sort.Sort(sort.Reverse(keys))
	return keys
}
