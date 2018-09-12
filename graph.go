package main

import (
	"fmt"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"time"
)

type Node struct {
	Id    int64
	Name  string
	IsDir bool
}

type TotalInfo struct {
	ReadingTime    time.Duration
	CountFiles     int64
	CountFolders   int64
	TotalFilesSize uint64
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

func createFileSystemGraph(path string) (graph *simple.WeightedDirectedGraph, root *Node, total TotalInfo) {
	total = TotalInfo{}
	graph = simple.NewWeightedDirectedGraph(0, 0)

	start := time.Now()

	var nodeid int64 = 0
	root = &Node{Id: nodeid, Name: path, IsDir: true}
	graph.AddNode(root)
	nodeid++

	queue := make([]*DirNode, 0)
	queue = append(queue, &DirNode{Node: root, Path: path})

	for len(queue) > 0 {
		curr := queue[0]
		parent := curr.Node

		for _, entry := range dirents(curr.Path) {
			fullPath := filepath.Join(curr.Path, entry.Name())

			node := &Node{Id: nodeid, Name: entry.Name(), IsDir: entry.IsDir()}
			graph.AddNode(node)
			nodeid++

			if entry.IsDir() {
				queue = append(queue, &DirNode{Node: node, Path: fullPath})
				edge := graph.NewWeightedEdge(parent, node, 0)
				graph.SetWeightedEdge(edge)
				total.CountFolders++
			} else {
				total.CountFiles++
				sz := uint64(entry.Size())
				total.TotalFilesSize += sz

				edge := graph.NewWeightedEdge(parent, node, float64(sz))
				graph.SetWeightedEdge(edge)
			}
		}

		queue = queue[1:]
	}
	total.ReadingTime = time.Since(start)
	return graph, root, total
}
