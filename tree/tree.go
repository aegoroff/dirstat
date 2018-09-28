// Package github is a Red-black search binary tree implementation
package tree

const (
    Black = iota
    Red
)

type RbTree struct {
    Root *Node
    tnil *Node
}

// Represent red-black tree node
type Node struct {
    // Node key (data)
    Key    *Comparable
    Parent *Node
    Left   *Node
    Right  *Node
    // Node color (red or black)
    Color int
    // Subtree size including node itself
    Size int64
}

type Comparable interface {
    LessThan(y interface{}) bool
    EqualTo(y interface{}) bool
}

// Creates new Red-Black empty tree
func NewRbTree() *RbTree {
    tnil := Node{Color: Black}
    return &RbTree{tnil: &tnil}
}

// Creates new node
func NewNode(si Comparable) *Node {
    return &Node{Key: &si}
}

// Walks tree inorder (left, node, right)
func WalkInorder(root *Node, action func(*Node)) {
    if root != nil && root.Key != nil {
        WalkInorder(root.Left, action)
        action(root)
        WalkInorder(root.Right, action)
    }
}

// Walks tree preorder (node, left, right)
func WalkPreorder(root *Node, action func(*Node)) {
    if root != nil && root.Key != nil {
        action(root)
        WalkPreorder(root.Left, action)
        WalkPreorder(root.Right, action)
    }
}
