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
	"strings"
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

type TotalInfo struct {
	ReadingTime    time.Duration
	CountFiles     int64
	CountFolders   int64
	TotalFilesSize uint64
}

var pathSeparator = fmt.Sprintf("%c", os.PathSeparator)

func main() {
	options := Options{}

	goptions.ParseAndFail(&options)

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist. Details:\n  %v", options.Path, err)
	}

	fmt.Printf("Root: %s\n\n", options.Path)

	gr, root, elapsed := createFileSystemGraph(options.Path)

	analyzeGraphAndOutputResults(gr, root, elapsed, options)

	printMemUsage()
}

func analyzeGraphAndOutputResults(gr *simple.WeightedDirectedGraph, root *Node, elapsed time.Duration, options Options) {
	verboseRanges := make(map[int]bool)
	for _, x := range options.Range {
		verboseRanges[x] = true
	}

	total := TotalInfo{ReadingTime: elapsed}
	allPaths := path.DijkstraFrom(root, gr)
	stat := make(map[Range]int64)
	bfs := traverse.BreadthFirst{}

	fileNodesByRange := map[Range][]*Node{}

	includeIntoRange := func(i int) bool {
		return options.Verbosity && verboseRanges[i+1]
	}

	bfs.Walk(gr, root, func(n graph.Node, d int) bool {
		nn := n.(*Node)
		if nn.IsDir {
			total.CountFolders++
		} else {
			_, w := allPaths.To(nn.Id)
			total.CountFiles++
			total.TotalFilesSize += uint64(w)
			for i, r := range fileSizeRanges {
				if !r.contains(w) {
					continue
				}
				stat[r]++

				if includeIntoRange(i) {
					nodes, ok := fileNodesByRange[r]
					if !ok {
						fileNodesByRange[r] = []*Node{nn}
					} else {
						fileNodesByRange[r] = append(nodes, nn)
					}
				}
			}
		}

		return false
	})

	fmt.Printf("Total files stat:\n")

	const format = "%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

	fmt.Fprintf(tw, format, "File size", "Amount", "Percent")
	fmt.Fprintf(tw, format, "---------", "------", "-------")

	var heads []string
	for _, r := range fileSizeRanges {
		percent := (float64(stat[r]) / float64(total.CountFiles)) * 100
		head := fmt.Sprintf("Between %s and %s", humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)))
		heads = append(heads, head)

		fmt.Fprintf(tw, "%v\t%v\t%.2f%%\n", head, stat[r], percent)
	}
	tw.Flush()

	if options.Verbosity && len(options.Range) > 0 {
		fmt.Printf("\nDetailed files stat:\n")
		for i, r := range fileSizeRanges {
			if includeIntoRange(i) && stat[r] > 0 {
				fmt.Printf("%s\n", heads[i])
				outputFilesInfoWithinRange(fileNodesByRange[r], &allPaths, r)
			}
		}
	}

	printTotals(total)
}

func outputFilesInfoWithinRange(nodes []*Node, allPaths *path.Shortest, r Range) {
	for _, node := range nodes {
		if node.IsDir {
			continue
		}

		nodes, w := allPaths.To(node.Id)

		if !r.contains(w) {
			continue
		}

		fullPath := makeFullPath(nodes)

		fmt.Printf("   %s - %s\n", fullPath, humanize.IBytes(uint64(w)))
	}
}

func makeFullPath(nodes []graph.Node) string {
	var parts []string
	for _, p := range nodes {
		n := p.(*Node).Name

		if strings.LastIndex(n, pathSeparator) == len(n)-1 {
			n = strings.TrimRight(n, pathSeparator)
		}

		parts = append(parts, n)
	}
	fullPath := strings.Join(parts, pathSeparator)
	return fullPath
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
