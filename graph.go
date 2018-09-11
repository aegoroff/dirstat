package main

import "fmt"

type Node struct {
	Id    int64
	Name  string
	IsDir bool
}

type FileItem struct {
	Node *Node
	Path string
}

func (n Node) ID() int64 {
	return n.Id
}

func (n Node) DOTID() string {
	return fmt.Sprintf("\"%s\"", n.Name)
}
