package main

import (
	"fmt"
	"github.com/spf13/afero"
	"math"
	"path/filepath"
	"time"

	"gonum.org/v1/gonum/graph/simple"
)

type node struct {
	NodeID int64
	Name   string
	IsDir  bool
}

type dirNode struct {
	Node *node
	Path string
}

type walkNode struct {
	Node   *node
	Parent string
	Size   int64
}

func (n node) ID() int64 {
	return n.NodeID
}

func (n node) DOTID() string {
	return fmt.Sprintf("\"%s\"", n.Name)
}

func createFileSystemGraph(path string, fs afero.Fs) (graph *simple.WeightedDirectedGraph, root *node, elapsed time.Duration) {
	graph = simple.NewWeightedDirectedGraph(0, math.Inf(1))

	start := time.Now()

	root = &node{NodeID: 0, Name: path, IsDir: true}
	graph.AddNode(root)

	ch := make(chan *walkNode, 1024)

	go runWalkingDir(path, fs, 1, ch)

	queue := []*dirNode{{Node: root, Path: path}}

	for walkNode := range ch {
		node := walkNode.Node

		graph.AddNode(node)

		for queue[0].Path != walkNode.Parent {
			queue = queue[1:]
		}

		parentNode := queue[0]

		if node.IsDir {
			fullPath := filepath.Join(walkNode.Parent, node.Name)
			queue = append(queue, &dirNode{Node: node, Path: fullPath})
		}

		weight := float64(walkNode.Size)
		edge := graph.NewWeightedEdge(parentNode.Node, node, weight)
		graph.SetWeightedEdge(edge)
	}

	elapsed = time.Since(start)
	return graph, root, elapsed
}

func runWalkingDir(path string, fs afero.Fs, nextID int64, ch chan<- *walkNode) {
	filesystemCh := make(chan filesystemItem, 1024)
	go func() {
		walkDirBreadthFirst(path, fs, filesystemCh)
	}()

	go func() {
		defer close(ch)
		for item := range filesystemCh {
			if item.event == fsEventDir {
				// Continue
			} else {
				node := &node{NodeID: nextID, Name: item.entry.Name(), IsDir: item.entry.IsDir()}
				ch <- &walkNode{Node: node, Parent: item.dir, Size: item.entry.Size()}
				nextID++
			}
		}
	}()
}
