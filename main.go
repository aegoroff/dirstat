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
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Options struct {
	Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
	Verbosity bool          `goptions:"-v, --verbose, description='Be verbose'"`
	Range     []int         `goptions:"-r, --range, description='Output verbose files info for ranges specified'"`
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

const totalTemplate = `
Read taken:    {{.ReadingTime}}
Sort taken:    {{.SortingTime}}

Total files:   {{.CountFiles}} ({{.TotalFilesSize | toBytesString }})
Total folders: {{.CountFolders}}
`

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

	start = time.Now()
	sorted, _ := topo.Sort(fileSystemGraph)
	totalStat.SortingTime = time.Since(start)

	allPaths := path.DijkstraFrom(sorted[0], fileSystemGraph)

	ranges := [...]Range{
		{Min: 0, Max: 100 * Kbyte},
		{Min: 100 * Kbyte, Max: Mbyte},
		{Min: Mbyte, Max: 10 * Mbyte},
		{Min: 10 * Mbyte, Max: 100 * Mbyte},
		{Min: 100 * Mbyte, Max: Gbyte},
		{Min: Gbyte, Max: 10 * Gbyte},
		{Min: 10 * Gbyte, Max: 100 * Gbyte},
		{Min: 100 * Gbyte, Max: Tbyte},
	}

	stat := make(map[Range]int64)

	bfs := traverse.BreadthFirst{}
	bfs.Walk(fileSystemGraph, sorted[0], func(n graph.Node, d int) bool {
		nn := n.(*Node)
		if !nn.IsDir {
			_, w := allPaths.To(nn.Id)

			sz := int64(w)
			for _, r := range ranges {
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

	for i, r := range ranges {
		fmt.Printf("Total files between %s and %s: %d\n", humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)), stat[r])
		if options.Verbosity && verboseRanges[i+1] {
			outputFilesInfoWithinRange(sorted, &allPaths, r.Min, r.Max)
		}
	}

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
	report.Execute(os.Stdout, totalStat)

	PrintMemUsage()
}

func outputFilesInfoWithinRange(sorted []graph.Node, allPaths *path.Shortest, min int64, max int64) {
	var filesCount uint64
	for _, node := range sorted {
		n := node.(*Node)
		if n.IsDir {
			continue
		}
		paths, w := allPaths.To(n.Id)

		if w < float64(min) || w > float64(max) {
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
