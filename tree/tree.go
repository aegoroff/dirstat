package tree

const (
    Black = iota
    Red
)

type RbTree struct {
    Root *TreeNode
    tnil *TreeNode
}

type TreeNode struct {
    Key    *Comparable
    Parent *TreeNode
    Left   *TreeNode
    Right  *TreeNode
    Color  int
    Size   int64
}

type Comparable interface {
    LessThan(y interface{}) bool
    EqualTo(y interface{}) bool
}

// Creates new Red-Black empty tree
func NewRbTree() *RbTree {
    tnil := TreeNode{Color: Black}
    return &RbTree{tnil: &tnil}
}

func TreeWalkInorder(root *TreeNode, action func(*TreeNode)) {
    if root != nil && root.Key != nil {
        TreeWalkInorder(root.Left, action)
        action(root)
        TreeWalkInorder(root.Right, action)
    }
}

func TreeWalkPreorder(root *TreeNode, action func(*TreeNode)) {
    if root != nil && root.Key != nil {
        action(root)
        TreeWalkPreorder(root.Left, action)
        TreeWalkPreorder(root.Right, action)
    }
}

func TreeInsert(root *TreeNode, z *TreeNode) {
    var y *TreeNode
    x := root
    for x != nil {
        y = x
        if (*z.Key).LessThan(*x.Key) {
            x = x.Left
        } else {
            x = x.Right
        }
    }

    z.Parent = y
    if y == nil {
        root = z
    } else if (*z.Key).LessThan(*y.Key) {
        y.Left = z
    } else {
        y.Right = z
    }
}

// Inserts new node into Red-Black tree
// Creates Root if tree is empty
func RbTreeInsert(tree *RbTree, z *TreeNode) {
    if tree.Root == nil {
        tree.Root = z
        tree.Root.Color = Black
        tree.Root.Parent = tree.tnil
        tree.Root.Left = tree.tnil
        tree.Root.Right = tree.tnil
        tree.Root.Size = 1
        return
    }
    y := tree.tnil
    x := tree.Root
    z.Size = 1
    for x != tree.tnil {
        y = x
        y.Size++
        if (*z.Key).LessThan(*x.Key) {
            x = x.Left
        } else {
            x = x.Right
        }
    }

    z.Parent = y
    if y == tree.tnil {
        tree.Root = z
    } else if (*z.Key).LessThan(*y.Key) {
        y.Left = z
    } else {
        y.Right = z
    }
    z.Left = tree.tnil
    z.Right = tree.tnil
    z.Color = Red
    rbInsertFixup(tree, z)
}

func rbInsertFixup(tree *RbTree, z *TreeNode) {
    for z.Parent.Color == Red {
        if z.Parent == z.Parent.Parent.Left {
            y := z.Parent.Parent.Right
            if y.Color == Red {
                z.Parent.Color = Black
                y.Color = Black
                z.Parent.Parent.Color = Red
                z = z.Parent.Parent
            } else if z == z.Parent.Right {
                z = z.Parent
                leftRotate(tree, z)
            } else {
                z.Parent.Color = Black
                z.Parent.Parent.Color = Red
                rightRotate(tree, z.Parent.Parent)
            }
        } else {
            y := z.Parent.Parent.Left
            if y.Color == Red {
                z.Parent.Color = Black
                y.Color = Black
                z.Parent.Parent.Color = Red
                z = z.Parent.Parent
            } else if z == z.Parent.Left {
                z = z.Parent
                rightRotate(tree, z)
            } else {
                z.Parent.Color = Black
                z.Parent.Parent.Color = Red
                leftRotate(tree, z.Parent.Parent)
            }
        }
    }
    tree.Root.Color = Black
}

// Searches value specified within search tree
func TreeSearch(root *TreeNode, value *Comparable) (*TreeNode, bool) {
    var x *TreeNode
    x = root
    for x != nil && x.Key != nil && !(*value).EqualTo(*x.Key) {
        if (*value).LessThan(*x.Key) {
            x = x.Left
        } else {
            x = x.Right
        }
    }
    return x, x != nil && x.Key != nil
}

func TreeMinimum(root *TreeNode) *TreeNode {
    x := root
    for x.Left != nil && x.Left.Key != nil {
        x = x.Left
    }
    return x
}

func TreeMaximum(root *TreeNode) *TreeNode {
    x := root
    for x.Right != nil && x.Right.Key != nil {
        x = x.Right
    }
    return x
}

func TreeSuccessor(n *TreeNode) *TreeNode {
    if n.Right != nil && n.Right.Key != nil {
        return TreeMinimum(n.Right)
    }

    y := n.Parent
    for y != nil && y.Key != nil && n == y.Right {
        n = y
        y = y.Parent
    }

    return y
}

func TreePredecessor(n *TreeNode) *TreeNode {
    if n.Left != nil && n.Left.Key != nil {
        return TreeMaximum(n.Left)
    }

    y := n.Parent
    for y != nil && y.Key != nil && n == y.Left {
        n = y
        y = y.Parent
    }

    return y
}

// Deletes node specified from binary search tree
func TreeDelete(root *TreeNode, z *TreeNode) {
    if z.Left == nil {
        transplant(root, z, z.Right)
    } else if z.Right == nil {
        transplant(root, z, z.Left)
    } else {
        y := TreeMinimum(z.Right)
        if y.Parent != z {
            transplant(root, y, y.Right)
            y.Right = z.Right
            y.Right.Parent = y
        }
        transplant(root, z, y)
        y.Left = z.Left
        y.Left.Parent = y
    }
}

func OrderStatisticSelect(root *TreeNode, i int64) *TreeNode {
    r := root.Left.Size + 1
    if i == r {
        return root
    } else if i < r {
        return OrderStatisticSelect(root.Left, i)
    } else {
        return OrderStatisticSelect(root.Right, i - r)
    }
}

// Deletes node specified from Red-black tree
func RbTreeDelete(tree *RbTree, z *TreeNode) {
    y := z

    p := z.Parent
    for p != tree.tnil {
        p.Size--
        p = p.Parent
    }

    var x *TreeNode
    yOriginalColor := y.Color
    if z.Left == tree.tnil {
        x = z.Right
        rbTransplant(tree, z, z.Right)
    } else if z.Right == tree.tnil {
        x = z.Left
        rbTransplant(tree, z, z.Left)
    } else {
        y := TreeMinimum(z.Right)
        yOriginalColor = y.Color
        x = y.Right
        if y.Parent == z {
            x.Parent = y
        } else {
            rbTransplant(tree, y, y.Right)
            y.Right = z.Right
            y.Right.Parent = y
        }
        rbTransplant(tree, z, y)
        y.Left = z.Left
        y.Left.Parent = y
        y.Color = z.Color
    }
    if yOriginalColor == Black {
        rbDeleteFixup(tree, x)
    }
}

func rbDeleteFixup(tree *RbTree, x *TreeNode) {
    for x != tree.Root && x.Color == Black {
        if x == x.Parent.Left {
            w := x.Parent.Right
            if w.Color == Red {
                w.Color = Black
                x.Parent.Color = Red
                leftRotate(tree, x.Parent)
                w = x.Parent.Right
            }
            if w.Left.Color == Black && w.Right.Color == Black {
                w.Color = Red
                x = x.Parent
            } else if w.Right.Color == Black {
                w.Left.Color = Black
                w.Color = Red
                rightRotate(tree, w)
                w = x.Parent.Right
            } else {
                w.Color = x.Parent.Color
                x.Parent.Color = Black
                w.Right.Color = Black
                leftRotate(tree, x.Parent)
                x = tree.Root
            }
        } else {
            w := x.Parent.Left
            if w.Color == Red {
                w.Color = Black
                x.Parent.Color = Red
                rightRotate(tree, x.Parent)
                w = x.Parent.Left
            }
            if w.Right.Color == Black && w.Left.Color == Black {
                w.Color = Red
                x = x.Parent
            } else if w.Left.Color == Black {
                w.Right.Color = Black
                w.Color = Red
                leftRotate(tree, w)
                w = x.Parent.Left
            } else {
                w.Color = x.Parent.Color
                x.Parent.Color = Black
                w.Left.Color = Black
                rightRotate(tree, x.Parent)
                x = tree.Root
            }
        }
    }
}

func transplant(root *TreeNode, u *TreeNode, v *TreeNode) {
    if u.Parent == nil {
        root = v
    } else if u == u.Parent.Left {
        u.Parent.Left = v
    } else {
        u.Parent.Right = v
    }
    if v != nil {
        v.Parent = u.Parent
    }
}

func rbTransplant(tree *RbTree, u *TreeNode, v *TreeNode) {
    if u.Parent == tree.tnil {
        tree.Root = v
    } else if u == u.Parent.Left {
        u.Parent.Left = v
    } else {
        u.Parent.Right = v
    }
    v.Parent = u.Parent
}

func leftRotate(tree *RbTree, x *TreeNode) {
    y := x.Right
    x.Right = y.Left
    if y.Left != tree.tnil {
        y.Left.Parent = x
    }
    y.Parent = x.Parent
    if x.Parent == tree.tnil {
        tree.Root = y
    } else if x == x.Parent.Left {
        x.Parent.Left = y
    } else {
        x.Parent.Right = y
    }

    y.Left = x
    x.Parent = y

    y.Size = x.Size
    x.Size = x.Left.Size + x.Right.Size + 1
}

func rightRotate(tree *RbTree, x *TreeNode) {
    y := x.Left
    x.Left = y.Right
    if y.Right != tree.tnil {
        y.Right.Parent = x
    }
    y.Parent = x.Parent
    if x.Parent == tree.tnil {
        tree.Root = y
    } else if x == x.Parent.Right {
        x.Parent.Right = y
    } else {
        x.Parent.Left = y
    }

    y.Right = x
    x.Parent = y

    y.Size = x.Size
    x.Size = x.Left.Size + x.Right.Size + 1
}
