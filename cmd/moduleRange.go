package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"path/filepath"
	"sort"
	"text/tabwriter"
)

type moduleRange struct {
	distribution  map[Range]containers
	aggregate     map[Range]fileStat
	verbose       bool
	enabledRanges map[int]bool
}

func (m *moduleRange) postScan() {

}

func (m *moduleRange) handler() fileHandler {
	return func(f *fileEntry) {
		unsignedSize := uint64(f.Size)

		// Calculate files range statistic
		for i, r := range fileSizeRanges {
			if !r.contains(f.Size) {
				continue
			}

			s := m.aggregate[r]
			s.TotalFilesCount++
			s.TotalFilesSize += unsignedSize
			m.aggregate[r] = s

			// Store each file info within range only i verbose option set
			if !m.verbose || !m.enabledRanges[i+1] {
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
}

func (m *moduleRange) output(tw *tabwriter.Writer, w io.Writer) {
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
