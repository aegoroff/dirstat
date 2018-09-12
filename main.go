// This tool shows directory specified statistic.  This includes files and dirs count, size etc.
package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/golang-collections/collections/stack"
	"github.com/voxelbrain/goptions"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"gonum.org/v1/gonum/graph/traverse"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type Options struct {
	Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
	Verbosity bool          `goptions:"-v, --verbose, description='Be verbose'"`
	Range     []int         `goptions:"-r, --range, description='Output verbose files info for fileSizeRanges specified'"`
	Path      string        `goptions:"-p, --path, obligatory, description='Name to the directory'"`
}

type Range struct {
	Min int64
	Max int64
}

const (
	_ int64 = 1 << (10 * iota)
	Kbyte
	Mbyte
	Gbyte
	Tbyte
)

type TotalStat struct {
	ReadingTime    time.Duration
	SortingTime    time.Duration
	CountFiles     int64
	CountFolders   int64
	TotalFilesSize uint64
}

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

func main() {
	options := Options{}

	goptions.ParseAndFail(&options)

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist. Details:\n  %v", options.Path, err)
	}

	fmt.Printf("Root: %s\n\n", options.Path)

	totalStat := TotalStat{}

	fileSystemGraph := simple.NewWeightedDirectedGraph(0, 0)

	walkingStack := stack.New()

	var nodeid int64 = 0

	start := time.Now()

	filepath.Walk(options.Path, func(path string, info os.FileInfo, err error) error {

		node := &Node{Id: nodeid, Name: info.Name(), IsDir: info.IsDir()}
		fileSystemGraph.AddNode(node)
		nodeid++

		var parent *FileItem

		if walkingStack.Len() > 0 {
			parent = walkingStack.Peek().(*FileItem)
			for !strings.HasPrefix(path, parent.Path) && walkingStack.Len() > 0 {
				parent = walkingStack.Pop().(*FileItem)
				if strings.HasPrefix(path, parent.Path) {
					walkingStack.Push(parent)
					break
				}
			}
		}

		if info.IsDir() {
			if parent != nil {
				edge := fileSystemGraph.NewWeightedEdge(parent.Node, node, 0)
				fileSystemGraph.SetWeightedEdge(edge)
			}

			walkingStack.Push(&FileItem{Path: path, Node: node})
			totalStat.CountFolders++
		} else {

			totalStat.CountFiles++
			sz := uint64(info.Size())
			totalStat.TotalFilesSize += sz

			edge := fileSystemGraph.NewWeightedEdge(parent.Node, node, float64(sz))
			fileSystemGraph.SetWeightedEdge(edge)
		}
		return nil
	})

	totalStat.ReadingTime = time.Since(start)

	printStatistic(fileSystemGraph, options, totalStat)

	printTotals(totalStat)

	printMemUsage()
}

func printStatistic(fileSystemGraph *simple.WeightedDirectedGraph, options Options, totalStat TotalStat) {
	start := time.Now()
	sorted, _ := topo.Sort(fileSystemGraph)
	totalStat.SortingTime = time.Since(start)

	allPaths := path.DijkstraFrom(sorted[0], fileSystemGraph)
	stat := make(map[Range]int64)
	bfs := traverse.BreadthFirst{}
	bfs.Walk(fileSystemGraph, sorted[0], func(n graph.Node, d int) bool {
		nn := n.(*Node)
		if !nn.IsDir {
			_, w := allPaths.To(nn.Id)

			sz := int64(w)
			for _, r := range fileSizeRanges {
				if sz < r.Min || sz > r.Max {
					continue
				}
				stat[r]++
			}

		}

		return false
	})
	verboseRanges := make(map[int]bool)
	for _, x := range options.Range {
		verboseRanges[x] = true
	}

	fmt.Printf("Total files stat:\n")
	fmt.Printf("%-35s%-12s %s\n", "", "Amount", "Percent")
	var heads []string
	for _, r := range fileSizeRanges {
		percent := (float64(stat[r]) / float64(totalStat.CountFiles)) * 100
		head := fmt.Sprintf("Between %s and %s", humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)))
		heads = append(heads, head)
		fmt.Printf("%-35s%-12d %.2f%%\n", head, stat[r], percent)
	}

	if options.Verbosity && len(options.Range) > 0 {
		fmt.Printf("\nDetailed files stat:\n")
		for i, r := range fileSizeRanges {
			if options.Verbosity && verboseRanges[i+1] {
				fmt.Printf("%s\n", heads[i])
				outputFilesInfoWithinRange(sorted, &allPaths, r)
			}
		}
	}
}

func printTotals(t TotalStat) {

	const totalTemplate = `
Read taken:    {{.ReadingTime}}
Sort taken:    {{.SortingTime}}

Total files:   {{.CountFiles}} ({{.TotalFilesSize | toBytesString }})
Total folders: {{.CountFolders}}
`

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
	report.Execute(os.Stdout, t)
}

func outputFilesInfoWithinRange(sorted []graph.Node, allPaths *path.Shortest, r Range) {
	var filesCount uint64
	for _, node := range sorted {
		n := node.(*Node)
		if n.IsDir {
			continue
		}
		paths, w := allPaths.To(n.Id)

		if w < float64(r.Min) || w > float64(r.Max) {
			continue
		}

		filesCount++

		var parts []string
		for _, p := range paths {
			n := p.(*Node).Name
			if n == "\\" {
				n = ""
			}

			parts = append(parts, n)
		}
		fullPath := strings.Join(parts, "/")

		fmt.Printf("   %s %s\n", fullPath, humanize.IBytes(uint64(w)))
	}
}
