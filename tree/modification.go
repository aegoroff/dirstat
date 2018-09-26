// This file contains all RB tree modification methods implementations
package tree

// Inserts new node into Red-Black tree
// Creates Root if tree is empty
func Insert(tree *RbTree, z *Node) {
    if z == nil {
        return
    }

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

func rbInsertFixup(tree *RbTree, z *Node) {
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

// Deletes node specified from Red-black tree
func Delete(tree *RbTree, z *Node) {
    if z == nil {
        return
    }

    y := z

    p := z.Parent
    for p != tree.tnil {
        p.Size--
        p = p.Parent
    }

    var x *Node
    yOriginalColor := y.Color
    if z.Left == tree.tnil {
        x = z.Right
        rbTransplant(tree, z, z.Right)
    } else if z.Right == tree.tnil {
        x = z.Left
        rbTransplant(tree, z, z.Left)
    } else {
        y := Minimum(z.Right)
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

func rbDeleteFixup(tree *RbTree, x *Node) {
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
            } else {
                if w.Right.Color == Black {
                    w.Left.Color = Black
                    w.Color = Red
                    rightRotate(tree, w)
                    w = x.Parent.Right
                }

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
            } else {
                if w.Left.Color == Black {
                    w.Right.Color = Black
                    w.Color = Red
                    leftRotate(tree, w)
                    w = x.Parent.Left
                }

                w.Color = x.Parent.Color
                x.Parent.Color = Black
                w.Left.Color = Black
                rightRotate(tree, x.Parent)
                x = tree.Root
            }
        }
    }
    x.Color = Black
}

func rbTransplant(tree *RbTree, u *Node, v *Node) {
    if u.Parent == tree.tnil {
        tree.Root = v
    } else if u == u.Parent.Left {
        u.Parent.Left = v
    } else {
        u.Parent.Right = v
    }
    v.Parent = u.Parent
}

func leftRotate(tree *RbTree, x *Node) {
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

func rightRotate(tree *RbTree, x *Node) {
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
