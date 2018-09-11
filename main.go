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
	"time"
)

type Options struct {
	Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
	Verbosity bool          `goptions:"-v, --verbose, description='Be verbose'"`
	Path      string        `goptions:"-p, --path, obligatory, description='Name to the directory'"`
}

type Node struct {
	Id    int64
	Name  string
	IsDir bool
}

type FileItem struct {
	Node Node
	Path string
}

func (n Node) ID() int64 {
	return n.Id
}

func (n Node) DOTID() string {
	return fmt.Sprintf("\"%s\"", n.Name)
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

func main() {
	options := Options{}

	goptions.ParseAndFail(&options)

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist. Details:\n  %v", options.Path, err)
	}

	fmt.Printf("Root: %s\n\n", options.Path)

	countFiles := 0
	countDirs := 0
	var totalSize uint64

	fileSystemGraph := simple.NewWeightedDirectedGraph(0, 0)

	walkingStack := stack.New()

	var nodeid int64 = 0

	start := time.Now()

	filepath.Walk(options.Path, func(path string, info os.FileInfo, err error) error {

		node := Node{Id: nodeid, Name: info.Name(), IsDir: info.IsDir()}
		fileSystemGraph.AddNode(node)
		nodeid++

		var parent FileItem

		if walkingStack.Len() > 0 {
			parent = walkingStack.Peek().(FileItem)
			for !strings.HasPrefix(path, parent.Path) && walkingStack.Len() > 0 {
				parent = walkingStack.Pop().(FileItem)
				if strings.HasPrefix(path, parent.Path) {
					walkingStack.Push(parent)
					break
				}
			}
		}

		if info.IsDir() {
			if len(parent.Path) > 0 {
				edge := fileSystemGraph.NewWeightedEdge(parent.Node, node, 0)
				fileSystemGraph.SetWeightedEdge(edge)
			}

			walkingStack.Push(FileItem{Path: path, Node: node})
			countDirs++
		} else {

			countFiles++
			sz := uint64(info.Size())
			totalSize += sz

			edge := fileSystemGraph.NewWeightedEdge(parent.Node, node, float64(sz))
			fileSystemGraph.SetWeightedEdge(edge)
		}
		return nil
	})

	readingTime := time.Since(start)

	start = time.Now()
	sorted, _ := topo.Sort(fileSystemGraph)
	sortingTime := time.Since(start)

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

	stat := [len(ranges)]int64{}

	bfs := traverse.BreadthFirst{}
	bfs.Walk(fileSystemGraph, sorted[0], func(n graph.Node, d int) bool {
		nn := n.(Node)
		if !nn.IsDir {
			_, w := allPaths.To(nn.Id)

			sz := int64(w)
			for i, r := range ranges {
				if sz < r.Min || sz > r.Max {
					continue
				}
				stat[i]++
			}

		}

		return false
	})

	for i, r := range ranges {
		fmt.Printf("Total files between %s and %s: %d\n", humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)), stat[i])
		if options.Verbosity {
			outputFilesInfoWithinRange(sorted, &allPaths, r.Min, r.Max)
		}
	}

	fmt.Printf("\nRead taken %v\n", readingTime)
	fmt.Printf("Sort taken %v\n\n", sortingTime)
	fmt.Printf("Total files %d Total size: %s\n", countFiles, humanize.IBytes(totalSize))
	fmt.Printf("Total folders %d\n", countDirs)
}

func outputFilesInfoWithinRange(sorted []graph.Node, allPaths *path.Shortest, min int64, max int64) {
	var filesCount uint64
	for _, node := range sorted {
		n := node.(Node)
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
			n := p.(Node).Name
			if n == "\\" {
				n = ""
			}

			parts = append(parts, n)
		}
		fullPath := strings.Join(parts, "/")

		fmt.Printf("   %s %s\n", fullPath, humanize.IBytes(uint64(w)))
	}
}
