// This tool shows directory specified statistic.  This includes files and dirs count, size etc.
package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/golang-collections/collections/stack"
	"github.com/voxelbrain/goptions"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Options struct {
	Help goptions.Help `goptions:"-h, --help, description='Show this help'"`
	Path string        `goptions:"-p, --path, obligatory, description='Name to the directory'"`
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

func main() {
	options := Options{}

	goptions.ParseAndFail(&options)

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist. Details:\n  %v", options.Path, err)
	}

	countFiles := 0
	countDirs := 0
	var totalSize uint64

	fileSystemGraph := simple.NewWeightedDirectedGraph(0, 0)

	walkingStack := stack.New()

	var nodeid int64 = 0

	start := time.Now()

	filepath.Walk(options.Path, func(path string, info os.FileInfo, err error) error {
		//fmt.Printf("%s", path)

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
			//fmt.Printf(" (%s)", humanize.Bytes(sz))
			totalSize += sz

			edge := fileSystemGraph.NewWeightedEdge(parent.Node, node, float64(sz))
			fileSystemGraph.SetWeightedEdge(edge)
		}
		//fmt.Print("\n")
		return nil
	})

	readingTime := time.Since(start)

	start = time.Now()
	sorted, _ := topo.Sort(fileSystemGraph)
	sortingTime := time.Since(start)

	//allp := path.DijkstraAllPaths(fileSystemGraph)

	root := sorted[0]

	allPaths := path.DijkstraFrom(root, fileSystemGraph)
	pathTo200, w200 := allPaths.To(433)

	var parts []string
	for _, p := range pathTo200 {
		n := p.(Node).Name
		if n == "\\" {
			n = ""
		}

		parts = append(parts, n)
	}

	p200s := strings.Join(parts, "/")

	fmt.Printf("Read taken %v\n", readingTime)
	fmt.Printf("Sort taken %v\n\n", sortingTime)
	fmt.Printf("Total files %d Total size: %s\n", countFiles, humanize.Bytes(totalSize))
	fmt.Printf("Total folders %d\n", countDirs)
	fmt.Printf("W200 (%s) %s\n", p200s, humanize.Bytes(uint64(w200)))
}
