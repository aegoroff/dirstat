package module

import (
	"dirstat/module/internal/sys"
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"path/filepath"
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

type moduleRange struct {
	distribution     map[Range]containers
	aggregate        map[Range]fileStat
	verbose          bool
	enabledRanges    []int
	enabledRangesMap map[int]bool
}

type moduleRangeNoOut struct {
	moduleRange
}

func (m *moduleRange) init() {
	m.enabledRangesMap = make(map[int]bool)
	for _, x := range m.enabledRanges {
		m.enabledRangesMap[x] = true
	}
}

func (m *moduleRange) postScan() {

}

func (m *moduleRange) folderHandler(_ *sys.FolderEntry) {

}

func (m *moduleRange) fileHandler(f *sys.FileEntry) {
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
		fileContainer := container{size: f.Size, name: filepath.Join(f.Parent, f.Name), count: 1}
		m.distribution[r] = append(nodes, &fileContainer)
	}
}

// Mute parent output
func (m *moduleRangeNoOut) output(*tabwriter.Writer, io.Writer) {

}

func (m *moduleRange) output(_ *tabwriter.Writer, w io.Writer) {
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
