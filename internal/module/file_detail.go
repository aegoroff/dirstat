package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/countingsort"
	"github.com/aegoroff/godatastruct/rbtree"
	"sort"
)

type detailsFile struct {
	enabled     map[int]*Range
	totalRanges int
}

type detailFileHandler struct {
	*detailsFile
}

type detailFileRenderer struct {
	*detailsFile
	*baseRenderer
	pd decorator
}

func newDetailsFile(ranges rbtree.RbTree, enabledRanges []int) *detailsFile {
	w := &detailsFile{
		enabled:     make(map[int]*Range, len(enabledRanges)),
		totalRanges: int(ranges.Len()),
	}

	for _, er := range enabledRanges {
		n, ok := ranges.OrderStatisticSelect(int64(er))
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

// Handler method

func (m *detailFileHandler) Handle(evt *scan.Event) {
	f := evt.File

	// Calculate files range statistic
	for _, r := range m.enabled {
		// Store each file info within range only if this file size detail option set
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

	for i := len(keys) - 1; i >= 0; i-- {
		k := keys[i]
		r := m.enabled[k]
		if len(r.files) == 0 {
			continue
		}

		p.Cprint("<gray>%2d. %s</>\n", k, r.String())

		sort.Sort(sort.Reverse(r.files))

		for _, f := range r.files {
			size := humanSize(f.size)
			p.Cprint("   %s - <yellow>%s</>\n", m.pd.decorate(f.String()), size)
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
	countingsort.Sort(keys, m.totalRanges)
	return keys
}
