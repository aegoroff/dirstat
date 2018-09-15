// This tool shows directory specified statistic.  This includes files and dirs count, size etc.
package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/voxelbrain/goptions"
	"log"
	"os"
	"path/filepath"
	"text/tabwriter"
	"text/template"
	"time"
)

type Options struct {
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

var fileSizeRanges = [...]Range{
	{Min: 0, Max: 100 * Kbyte},
	{Min: 100 * Kbyte, Max: Mbyte},
	{Min: Mbyte, Max: 10 * Mbyte},
	{Min: 10 * Mbyte, Max: 100 * Mbyte},
	{Min: 100 * Mbyte, Max: Gbyte},
	{Min: Gbyte, Max: 10 * Gbyte},
	{Min: 10 * Gbyte, Max: 100 * Gbyte},
	{Min: 100 * Gbyte, Max: Tbyte},
}

type FileStat struct {
	TotalFilesSize  uint64
	TotalFilesCount int64
}

type FileEntry struct {
	Size int64
	Path string
}

type WalkEntry struct {
	Size   int64
	Parent string
	Name   string
	IsDir  bool
}

type TotalInfo struct {
	ReadingTime    time.Duration
	CountFiles     int64
	CountFolders   int64
	TotalFilesSize uint64
}

func main() {
	options := Options{}

	goptions.ParseAndFail(&options)

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist. Details:\n  %v", options.Path, err)
	}

	fmt.Printf("Root: %s\n\n", options.Path)

	runAnalyze(options)

	printMemUsage()
}

func runAnalyze(options Options) {
	total, stat, fileNodesByRange := walk(options)

	fmt.Printf("Total files stat:\n")

	const format = "%v\t%v\t%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Fprintf(tw, format, "File size", "Amount", "%", "Size", "%")
	fmt.Fprintf(tw, format, "---------", "------", "------", "----", "------")

	var heads []string
	for i, r := range fileSizeRanges {
		percentOfCount := (float64(stat[r].TotalFilesCount) / float64(total.CountFiles)) * 100
		percentOfSize := (float64(stat[r].TotalFilesSize) / float64(total.TotalFilesSize)) * 100
		head := fmt.Sprintf("%d. Between %s and %s", i+1, humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)))
		heads = append(heads, head)

		fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", head, stat[r].TotalFilesCount, percentOfCount, humanize.IBytes(stat[r].TotalFilesSize), percentOfSize)
	}
	tw.Flush()

	if options.Verbosity && len(options.Range) > 0 {
		fmt.Printf("\nDetailed files stat:\n")
		for i, r := range fileSizeRanges {
			if len(fileNodesByRange[r]) == 0 {
				continue
			}

			fmt.Printf("%s\n", heads[i])
			for _, item := range fileNodesByRange[r] {
				fmt.Printf("   %s - %s\n", item.Path, humanize.IBytes(uint64(item.Size)))
			}
		}
	}

	printTotals(total)
}

func walk(options Options) (TotalInfo, map[Range]FileStat, map[Range][]*FileEntry) {
	verboseRanges := make(map[int]bool)
	for _, x := range options.Range {
		verboseRanges[x] = true
	}
	total := TotalInfo{}
	stat := make(map[Range]FileStat)
	fileNodesByRange := map[Range][]*FileEntry{}

	ch := make(chan *WalkEntry, 1024)

	start := time.Now()

	go func(ch chan<- *WalkEntry) {
		walkDirBreadthFirst(options.Path, func(parent string, entry os.FileInfo) {
			ch <- &WalkEntry{IsDir: entry.IsDir(), Size: entry.Size(), Parent: parent, Name: entry.Name()}
		})
		close(ch)
	}(ch)

	for {
		walkEntry, ok := <-ch
		if !ok {
			break
		}

		if walkEntry.IsDir {
			total.CountFolders++
		} else {
			// Accumulate file statistic
			sz := uint64(walkEntry.Size)
			total.CountFiles++
			total.TotalFilesSize += sz
			for i, r := range fileSizeRanges {
				if !r.contains(float64(sz)) {
					continue
				}

				s := stat[r]
				s.TotalFilesCount++
				s.TotalFilesSize += sz
				stat[r] = s

				if options.Verbosity && verboseRanges[i+1] {
					fullPath := filepath.Join(walkEntry.Parent, walkEntry.Name)
					nodes, ok := fileNodesByRange[r]
					if !ok {
						fileNodesByRange[r] = []*FileEntry{{Path: fullPath, Size: walkEntry.Size}}
					} else {
						fileNodesByRange[r] = append(nodes, &FileEntry{Path: fullPath, Size: walkEntry.Size})
					}
				}
			}
		}
	}

	total.ReadingTime = time.Since(start)
	return total, stat, fileNodesByRange
}

func printTotals(t TotalInfo) {

	const totalTemplate = `
Total files:   {{.CountFiles}} ({{.TotalFilesSize | toBytesString }})
Total folders: {{.CountFolders}}

Read taken:    {{.ReadingTime}}
`

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
	report.Execute(os.Stdout, t)
}
