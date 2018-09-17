package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"gonum.org/v1/gonum/graph/simple"
)

type Node struct {
	NodeID int64
	Name   string
	IsDir  bool
}

type DirNode struct {
	Node *Node
	Path string
}

type WalkNode struct {
	Node   *Node
	Parent string
	Size   int64
}

func (n Node) ID() int64 {
	return n.NodeID
}

func (n Node) DOTID() string {
	return fmt.Sprintf("\"%s\"", n.Name)
}

func createFileSystemGraph(path string) (graph *simple.WeightedDirectedGraph, root *Node, elapsed time.Duration) {
	graph = simple.NewWeightedDirectedGraph(0, math.Inf(1))

	start := time.Now()

	root = &Node{NodeID: 0, Name: path, IsDir: true}
	graph.AddNode(root)

	ch := make(chan *WalkNode, 1024)

	go runWalkingDir(path, 1, ch)

	queue := []*DirNode{{Node: root, Path: path}}

	for {
		walkNode, ok := <-ch
		if !ok {
			break
		}
		node := walkNode.Node

		graph.AddNode(node)

		for queue[0].Path != walkNode.Parent {
			queue = queue[1:]
		}

		parentNode := queue[0]

		if node.IsDir {
			fullPath := filepath.Join(walkNode.Parent, node.Name)
			queue = append(queue, &DirNode{Node: node, Path: fullPath})
		}

		weight := float64(walkNode.Size)
		edge := graph.NewWeightedEdge(parentNode.Node, node, weight)
		graph.SetWeightedEdge(edge)
	}

	elapsed = time.Since(start)
	return graph, root, elapsed
}

func runWalkingDir(path string, nextID int64, ch chan<- *WalkNode) {
	walkDirBreadthFirst(path, func(parent string, entry os.FileInfo) {
		node := &Node{NodeID: nextID, Name: entry.Name(), IsDir: entry.IsDir()}
		ch <- &WalkNode{Node: node, Parent: parent, Size: entry.Size()}
		nextID++
	})
	close(ch)
}
