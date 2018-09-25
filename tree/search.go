// This file contains all RB tree search methods implementations
package tree

// Searches value specified within search tree
func Search(root *Node, value *Comparable) (*Node, bool) {
    var x *Node
    x = root
    for x != nil && x.Key != nil && !(*value).EqualTo(*x.Key) {
        if (*value).LessThan(*x.Key) {
            x = x.Left
        } else {
            x = x.Right
        }
    }
    ok := x != nil && x.Key != nil

    if !ok {
        return nil, ok
    }

    return x, ok
}

// Gets tree's min element
func Minimum(root *Node) *Node {
    x := root
    for x.Left != nil && x.Left.Key != nil {
        x = x.Left
    }
    return x
}

// Gets tree's max element
func Maximum(root *Node) *Node {
    x := root
    for x.Right != nil && x.Right.Key != nil {
        x = x.Right
    }
    return x
}

// Gets node specified successor
func Successor(n *Node) *Node {
    if n.Right != nil && n.Right.Key != nil {
        return Minimum(n.Right)
    }

    y := n.Parent
    for y != nil && y.Key != nil && n == y.Right {
        n = y
        y = y.Parent
    }

    if y.Key == nil {
        return nil
    }

    return y
}

// Gets node specified predecessor
func Predecessor(n *Node) *Node {
    if n.Left != nil && n.Left.Key != nil {
        return Maximum(n.Left)
    }

    y := n.Parent
    for y != nil && y.Key != nil && n == y.Left {
        n = y
        y = y.Parent
    }

    if y.Key == nil {
        return nil
    }

    return y
}

// Gets i element from subtree
func OrderStatisticSelect(root *Node, i int64) *Node {
    r := root.Left.Size + 1
    if i == r {
        return root
    } else if i < r {
        return OrderStatisticSelect(root.Left, i)
    } else {
        return OrderStatisticSelect(root.Right, i-r)
    }
}
