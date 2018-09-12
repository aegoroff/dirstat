// This tool shows directory specified statistic.  This includes files and dirs count, size etc.
package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/voxelbrain/goptions"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
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

const (
	_ int64 = 1 << (10 * iota)
	Kbyte
	Mbyte
	Gbyte
	Tbyte
)

type TotalStat struct {
	ReadingTime    time.Duration
	CountFiles     int64
	CountFolders   int64
	TotalFilesSize uint64
}

type DirNode struct {
	Node *Node
	Path string
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

var pathSeparator = fmt.Sprintf("%c", os.PathSeparator)

func main() {
	options := Options{}

	goptions.ParseAndFail(&options)

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist. Details:\n  %v", options.Path, err)
	}

	fmt.Printf("Root: %s\n\n", options.Path)

	totalStat := TotalStat{}

	fileSystemGraph := simple.NewWeightedDirectedGraph(0, 0)

	var nodeid int64 = 0

	start := time.Now()

	root := &Node{Id: nodeid, Name: options.Path, IsDir: true}
	fileSystemGraph.AddNode(root)
	nodeid++

	queue := make([]*DirNode, 0)

	queue = append(queue, &DirNode{Node: root, Path: options.Path})

	for len(queue) > 0 {
		curr := queue[0]
		parent := curr.Node

		for _, entry := range dirents(curr.Path) {
			fullPath := filepath.Join(curr.Path, entry.Name())

			node := &Node{Id: nodeid, Name: entry.Name(), IsDir: entry.IsDir()}
			fileSystemGraph.AddNode(node)
			nodeid++

			if entry.IsDir() {
				queue = append(queue, &DirNode{Node: node, Path: fullPath})
				edge := fileSystemGraph.NewWeightedEdge(parent, node, 0)
				fileSystemGraph.SetWeightedEdge(edge)
				totalStat.CountFolders++
			} else {
				totalStat.CountFiles++
				sz := uint64(entry.Size())
				totalStat.TotalFilesSize += sz

				edge := fileSystemGraph.NewWeightedEdge(parent, node, float64(sz))
				fileSystemGraph.SetWeightedEdge(edge)
			}
		}

		queue = queue[1:]
	}

	totalStat.ReadingTime = time.Since(start)

	printStatistic(fileSystemGraph, root, options, totalStat)

	printTotals(totalStat)

	printMemUsage()
}

func printStatistic(fileSystemGraph *simple.WeightedDirectedGraph, root *Node, options Options, totalStat TotalStat) {

	allPaths := path.DijkstraFrom(root, fileSystemGraph)
	stat := make(map[Range]int64)
	bfs := traverse.BreadthFirst{}
	bfs.Walk(fileSystemGraph, root, func(n graph.Node, d int) bool {
		nn := n.(*Node)
		if !nn.IsDir {
			_, w := allPaths.To(nn.Id)

			for _, r := range fileSizeRanges {
				if !r.contains(w) {
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
			if options.Verbosity && verboseRanges[i+1] && stat[r] > 0 {
				fmt.Printf("%s\n", heads[i])
				outputFilesInfoWithinRange(fileSystemGraph.Nodes(), &allPaths, r)
			}
		}
	}
}

func printTotals(t TotalStat) {

	const totalTemplate = `
Total files:   {{.CountFiles}} ({{.TotalFilesSize | toBytesString }})
Total folders: {{.CountFolders}}

Read taken:    {{.ReadingTime}}
`

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
	report.Execute(os.Stdout, t)
}

func outputFilesInfoWithinRange(nodes []graph.Node, allPaths *path.Shortest, r Range) {
	for _, node := range nodes {
		n := node.(*Node)
		if n.IsDir {
			continue
		}
		paths, w := allPaths.To(n.Id)

		if !r.contains(w) {
			continue
		}

		var parts []string
		for _, p := range paths {
			n := p.(*Node).Name

			if strings.LastIndex(n, pathSeparator) == len(n)-1 {
				n = strings.TrimRight(n, pathSeparator)
			}

			parts = append(parts, n)
		}
		fullPath := strings.Join(parts, pathSeparator)

		fmt.Printf("   %s - %s\n", fullPath, humanize.IBytes(uint64(w)))
	}
}
