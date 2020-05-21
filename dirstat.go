// This tool shows directory specified statistic.  This includes files and dirs count, size etc.
package main

import (
	"fmt"
	"github.com/spf13/afero"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"github.com/voxelbrain/goptions"
)

type options struct {
	Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
	Verbosity bool          `goptions:"-v, --verbose, description='Be verbose'"`
	Range     []int         `goptions:"-r, --range, description='Output verbose files info for fileSizeRanges specified'"`
	Path      string        `goptions:"-p, --path, obligatory, description='Name to the directory'"`
}

const (
	_ int64 = 1 << (10 * iota)
	Kbyte
	Mbyte
	Gbyte
	Tbyte
)

const Top = 10

var fileSizeRanges = [...]Range{
	{Min: 0, Max: 100 * Kbyte},
	{Min: 100 * Kbyte, Max: Mbyte},
	{Min: Mbyte, Max: 10 * Mbyte},
	{Min: 10 * Mbyte, Max: 100 * Mbyte},
	{Min: 100 * Mbyte, Max: Gbyte},
	{Min: Gbyte, Max: 10 * Gbyte},
	{Min: 10 * Gbyte, Max: 100 * Gbyte},
	{Min: 100 * Gbyte, Max: Tbyte},
	{Min: Tbyte, Max: 10 * Tbyte},
	{Min: 10 * Tbyte, Max: 100 * Tbyte},
}

type fileStat struct {
	TotalFilesSize  uint64
	TotalFilesCount int64
}

type walkEntry struct {
	Size   int64
	Parent string
	Name   string
	IsDir  bool
}

type totalInfo struct {
	ReadingTime   time.Duration
	FilesTotal    countSizeAggregate
	CountFolders  int64
	CountFileExts int
}

type countSizeAggregate struct {
	Count int64
	Size  uint64
}

func main() {
	opt := options{}

	goptions.ParseAndFail(&opt)

	fs := afero.NewOsFs()

	if _, err := fs.Stat(opt.Path); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist. Details:\n  %v", opt.Path, err)
	}

	if opt.Path[len(opt.Path)-1] == ':' {
		opt.Path = filepath.Join(opt.Path, "\\")
	}

	w := os.Stdout

	_, _ = fmt.Fprintf(w, "Root: %s\n\n", opt.Path)

	runAnalyze(opt, fs, w)

	printMemUsage(w)
}

func runAnalyze(opt options, fs afero.Fs, w io.Writer) {
	total, stat, filesByRange, byExt, byFolder, topFiles := walk(opt, fs)
	total.CountFileExts = len(byExt)

	extBySize := createSliceFromMap(byExt, func(aggregate countSizeAggregate) int64 {
		return int64(aggregate.Size)
	})

	extByCount := createSliceFromMap(byExt, func(aggregate countSizeAggregate) int64 {
		return aggregate.Count
	})

	sort.Sort(sort.Reverse(extBySize))
	sort.Sort(sort.Reverse(extByCount))

	_, _ = fmt.Fprintf(w, "Total files stat:\n\n")

	const format = "%v\t%v\t%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	_, _ = fmt.Fprintf(tw, format, "File size", "Amount", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "------", "------", "----", "------")

	var heads []string
	for i, r := range fileSizeRanges {
		h := fmt.Sprintf("%d. Between %s and %s", i+1, humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)))
		heads = append(heads, h)

		count := stat[r].TotalFilesCount
		sz := stat[r].TotalFilesSize

		outputTopStatLine(tw, count, total, sz, h)
	}
	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d file extensions by size:\n\n", Top)
	_, _ = fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

	outputTopTenExtensions(tw, extBySize, total, func(data namedInts64, item *namedInt64) (int64, uint64) {
		count := byExt[item.name].Count
		sz := uint64(item.value)
		return count, sz
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d file extensions by count:\n\n", Top)
	_, _ = fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

	outputTopTenExtensions(tw, extByCount, total, func(data namedInts64, item *namedInt64) (int64, uint64) {
		count := item.value
		sz := byExt[item.name].Size
		return count, sz
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d files by size:\n\n", Top)
	_, _ = fmt.Fprintf(tw, "%v\t%v\n", "File", "Size")
	_, _ = fmt.Fprintf(tw, "%v\t%v\n", "------", "----")

	i := 1

	topFiles.Descend(func(c *rbtree.Comparable) bool {
		file := (*c).(namedInt64)
		h := fmt.Sprintf("%d. %s", i, file.name)

		i++

		sz := uint64(file.value)

		_, _ = fmt.Fprintf(tw, "%v\t%v\n", h, humanize.IBytes(sz))

		return true
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d folders by size:\n\n", Top)
	_, _ = fmt.Fprintf(tw, format, "Folder", "Files", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "------", "-----", "------", "----", "------")

	i = 1

	byFolder.Descend(func(c *rbtree.Comparable) bool {

		folder := (*c).(statItem)
		h := fmt.Sprintf("%d. %s", i, folder.name)

		i++

		count := folder.count
		sz := uint64(folder.size)

		outputTopStatLine(tw, count, total, sz, h)

		return true
	})

	_ = tw.Flush()

	if opt.Verbosity && len(opt.Range) > 0 {
		_, _ = fmt.Fprintf(w, "\nDetailed files stat:\n")
		for i, r := range fileSizeRanges {
			if len(filesByRange[r]) == 0 {
				continue
			}

			_, _ = fmt.Fprintf(w, "%s\n", heads[i])
			for _, item := range filesByRange[r] {
				fullPath := filepath.Join(item.Parent, item.Name)
				size := humanize.IBytes(uint64(item.Size))
				_, _ = fmt.Fprintf(w, "   %s - %s\n", fullPath, size)
			}
		}
	}

	printTotals(total, w)
}

func outputTopTenExtensions(tw *tabwriter.Writer, data namedInts64, total totalInfo, selector func(data namedInts64, item *namedInt64) (int64, uint64)) {
	for i := 0; i < Top && i < len(data); i++ {
		h := data[i].name

		count, sz := selector(data, data[i])

		outputTopStatLine(tw, count, total, sz, h)
	}
}

func outputTopStatLine(tw *tabwriter.Writer, count int64, total totalInfo, sz uint64, title string) {
	percentOfCount := countPercent(count, total)
	percentOfSize := sizePercent(sz, total)

	_, _ = fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
}

func countPercent(count int64, total totalInfo) float64 {
	return (float64(count) / float64(total.FilesTotal.Count)) * 100
}

func sizePercent(size uint64, total totalInfo) float64 {
	return (float64(size) / float64(total.FilesTotal.Size)) * 100
}

func createSliceFromMap(sizeByExt map[string]countSizeAggregate, mapper func(countSizeAggregate) int64) namedInts64 {
	var result = make(namedInts64, len(sizeByExt))
	i := 0
	for k, v := range sizeByExt {
		result[i] = &namedInt64{value: mapper(v), name: k}
		i++
	}
	return result
}

func walk(opt options, fs afero.Fs) (totalInfo, map[Range]fileStat, map[Range][]*walkEntry, map[string]countSizeAggregate, *rbtree.RbTree, *rbtree.RbTree) {
	verboseRanges := make(map[int]bool)
	for _, x := range opt.Range {
		verboseRanges[x] = true
	}
	total := totalInfo{}
	stat := make(map[Range]fileStat)
	filesByRange := make(map[Range][]*walkEntry)

	byExt := make(map[string]countSizeAggregate)

	ch := make(chan *walkEntry, 1024)

	start := time.Now()

	go func(ch chan<- *walkEntry) {
		defer close(ch)
		walkDirBreadthFirst(opt.Path, fs, func(parent string, entry os.FileInfo) {
			ch <- &walkEntry{IsDir: entry.IsDir(), Size: entry.Size(), Parent: parent, Name: entry.Name()}
		})
	}(ch)

	folderSizeTree := rbtree.NewRbTree()
	fileSizeTree := rbtree.NewRbTree()

	currFolderStat := statItem{}

	for we := range ch {
		if we.IsDir {
			total.CountFolders++
		} else {
			// Accumulate file statistic
			sz := uint64(we.Size)
			total.FilesTotal.Count++
			total.FilesTotal.Size += sz

			ext := filepath.Ext(we.Name)
			a := byExt[ext]
			a.Size += uint64(we.Size)
			a.Count++
			byExt[ext] = a

			if currFolderStat.name == "" {
				currFolderStat.name = we.Parent
			}

			if currFolderStat.name == we.Parent {
				currFolderStat.size += we.Size
				currFolderStat.count++
			} else {
				minfolder := folderSizeTree.Minimum()
				if folderSizeTree.Len() < Top || getSizeFromNode(minfolder) < currFolderStat.size {
					if folderSizeTree.Len() == Top {
						folderSizeTree.Delete(minfolder)
					}

					var c rbtree.Comparable
					c = currFolderStat
					node := rbtree.NewNode(&c)
					folderSizeTree.Insert(node)
				}

				currFolderStat.name = we.Parent
				currFolderStat.count = 1
				currFolderStat.size = we.Size
			}

			minfile := fileSizeTree.Minimum()
			if fileSizeTree.Len() < Top || getSizeFromNode(minfile) < we.Size {
				if fileSizeTree.Len() == Top {
					fileSizeTree.Delete(minfile)
				}

				fullPath := filepath.Join(we.Parent, we.Name)
				value := namedInt64{value: we.Size, name: fullPath}

				var c rbtree.Comparable
				c = value
				node := rbtree.NewNode(&c)
				fileSizeTree.Insert(node)
			}

			for i, r := range fileSizeRanges {
				if !r.contains(float64(sz)) {
					continue
				}

				s := stat[r]
				s.TotalFilesCount++
				s.TotalFilesSize += sz
				stat[r] = s

				if !opt.Verbosity || !verboseRanges[i+1] {
					continue
				}

				nodes, ok := filesByRange[r]
				if !ok {
					filesByRange[r] = []*walkEntry{we}
				} else {
					filesByRange[r] = append(nodes, we)
				}
			}
		}
	}

	total.ReadingTime = time.Since(start)
	return total, stat, filesByRange, byExt, folderSizeTree, fileSizeTree
}

func getSizeFromNode(node *rbtree.Node) int64 {
	if k, ok := (*node.Key).(statItem); ok {
		return k.size
	}

	if k, ok := (*node.Key).(namedInt64); ok {
		return k.value
	}

	return 0
}

func printTotals(t totalInfo, w io.Writer) {

	const totalTemplate = `
Total files:            {{.FilesTotal.Count}} ({{.FilesTotal.Size | toBytesString }})
Total folders:          {{.CountFolders}}
Total file extensions:  {{.CountFileExts}}

Read taken:    {{.ReadingTime}}
`

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
	_ = report.Execute(w, t)
}
