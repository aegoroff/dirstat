package main

import (
	"fmt"
	"gonum.org/v1/gonum/graph/simple"
	"math"
	"os"
	"path/filepath"
	"time"
)

type Node struct {
	Id    int64
	Name  string
	IsDir bool
}

type DirNode struct {
	Node *Node
	Path string
}

func (n Node) ID() int64 {
	return n.Id
}

func (n Node) DOTID() string {
	return fmt.Sprintf("\"%s\"", n.Name)
}

func createFileSystemGraph(path string) (graph *simple.WeightedDirectedGraph, root *Node, elapsed time.Duration) {
	graph = simple.NewWeightedDirectedGraph(0, math.Inf(1))

	start := time.Now()

	var nodeid int64 = 0
	root = &Node{Id: nodeid, Name: path, IsDir: true}
	graph.AddNode(root)
	nodeid++

	queue := []*DirNode{{Node: root, Path: path}}

	walkDirBreadthFirst(path, func(parent string, entry os.FileInfo) {
		node := &Node{Id: nodeid, Name: entry.Name(), IsDir: entry.IsDir()}
		graph.AddNode(node)
		nodeid++

		for queue[0].Path != parent {
			queue = queue[1:]
		}

		parentNode := queue[0]

		if entry.IsDir() {
			fullPath := filepath.Join(parent, node.Name)
			queue = append(queue, &DirNode{Node: node, Path: fullPath})
		}

		weight := float64(entry.Size())
		edge := graph.NewWeightedEdge(parentNode.Node, node, weight)
		graph.SetWeightedEdge(edge)
	})

	elapsed = time.Since(start)
	return graph, root, elapsed
}
