package tree

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "strings"
    "testing"
)

type IntNode struct {
    value int
}

type StringNode struct {
    value string
}

func (x IntNode) LessThan(y interface{}) bool {
    return x.value < (y.(IntNode)).value
}

func (x IntNode) EqualTo(y interface{}) bool {
    return x.value == (y.(IntNode)).value
}

func (x StringNode) LessThan(y interface{}) bool {
    return x.value < (y.(StringNode)).value
}

func (x StringNode) EqualTo(y interface{}) bool {
    return x.value == (y.(StringNode)).value
}

func TestInorderWalk(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    root := createIntegerTestTree()
    var result []int

    // Act
    WalkInorder(root, func(node *Node) {
        result = append(result, getIntValueOf(node))
    })

    // Assert
    ass.ElementsMatch(result, []int{2, 3, 4, 6, 7, 9, 13, 15, 17, 18, 20})
}

func TestInorderWalkStringTree(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    root := createStringTestTree()
    var result []string

    // Act
    WalkInorder(root, func(node *Node) {
        result = append(result, getStringValueOf(node))
    })

    // Assert
    ass.ElementsMatch(result, []string{"abc", "amd", "cisco", "do", "fake", "intel", "it", "let", "microsoft", "russia", "usa", "xxx", "yyy", "zen"})
}

func TestSearchSuccess(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    root := createIntegerTestTree()
    v := createIntNode(13)

    // Act
    found, ok := Search(root, v)

    // Assert
    ass.True(ok)
    ass.Equal(13, getIntValueOf(found))
}

func TestSearchFailure(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    root := createIntegerTestTree()
    v := createIntNode(22)

    // Act
    found, ok := Search(root, v)

    // Assert
    ass.False(ok)
    ass.Nil(found)
}

func TestSuccessor(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    root := createIntegerTestTree()
    v := createIntNode(13)
    r, _ := Search(root, v)

    // Act
    s := Successor(r)

    // Assert
    ass.Equal(15, getIntValueOf(s))
}

func TestPredecessor(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    root := createIntegerTestTree()

    // Act
    p := Predecessor(root)

    // Assert
    ass.Equal(13, getIntValueOf(p))
}

func TestDelete(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    root := createIntegerTestTree()
    v := createIntNode(6)
    r, _ := Search(root, v)
    var result []int

    // Act
    Delete(root, r)

    // Assert
    WalkInorder(root, func(node *Node) {
        result = append(result, getIntValueOf(node))
    })
    ass.ElementsMatch(result, []int{2, 3, 4, 7, 9, 13, 15, 17, 18, 20})
}

func TestTreeMinimum(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    root := createIntegerTestTree()

    // Act
    r := Minimum(root)

    // Assert
    ass.Equal(2, getIntValueOf(r))
}

func TestTreeMaximum(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    root := createIntegerTestTree()

    // Act
    r := Maximum(root)

    // Assert
    ass.Equal(20, getIntValueOf(r))
}

func TestRightRotate(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    r := Node{Key: createStringNode("root")}

    tree := NewRbTree()
    RbTreeInsert(tree, &r)

    y := Node{Key: createStringNode("y")}
    x := Node{Key: createStringNode("x")}
    a := Node{Key: createStringNode("a")}
    b := Node{Key: createStringNode("b")}
    g := Node{Key: createStringNode("g")}

    r.Right = &y
    y.Parent = &r
    y.Left = &x
    y.Right = &g
    x.Left = &a
    x.Right = &b
    x.Parent = &y
    g.Parent = &y
    a.Parent = &x
    b.Parent = &x

    // Act
    rightRotate(tree, &y)

    // Assert
    ass.Equal("root", getStringValueOf(x.Parent))
    ass.Equal("a", getStringValueOf(x.Left))
    ass.Equal("y", getStringValueOf(x.Right))
    ass.Equal("b", getStringValueOf(y.Left))
    ass.Equal("g", getStringValueOf(y.Right))
}

func TestLeftRotate(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    r := Node{Key: createStringNode("root")}

    tree := NewRbTree()
    RbTreeInsert(tree, &r)

    x := Node{Key: createStringNode("x")}
    y := Node{Key: createStringNode("y")}
    a := Node{Key: createStringNode("a")}
    b := Node{Key: createStringNode("b")}
    g := Node{Key: createStringNode("g")}

    r.Right = &x
    x.Parent = &r
    x.Left = &a
    x.Right = &y
    y.Left = &b
    y.Right = &g
    y.Parent = &y
    g.Parent = &y
    a.Parent = &x
    b.Parent = &y

    // Act
    leftRotate(tree, &x)

    // Assert
    ass.Equal("root", getStringValueOf(y.Parent))
    ass.Equal("x", getStringValueOf(y.Left))
    ass.Equal("g", getStringValueOf(y.Right))
    ass.Equal("a", getStringValueOf(x.Left))
    ass.Equal("b", getStringValueOf(x.Right))
}

func createIntegerTestTree() *Node {
    r := createIntNode(15)
    root := Node{Key: r}
    nodes := []int{6, 18, 3, 7, 2, 4, 13, 9, 17, 20}
    for _, n := range nodes {
        Insert(&root, &Node{Key: createIntNode(n)})
    }
    return &root
}

func TestRbTreeWalkInorder(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createTestRbTree()
    var result []string

    // Act
    WalkInorder(tree.Root, func(node *Node) {
        result = append(result, getStringValueOf(node))
    })

    // Assert
    ass.ElementsMatch(result, []string{"abc", "amd", "cisco", "do", "fake", "intel", "it", "let", "microsoft", "russia", "usa", "xxx", "yyy", "zen"})
}

func TestRbTreeWalkPreorder(t *testing.T) {
    // Arrange
    tree := createTestRbTree()
    b := strings.Builder{}

    // Act
    WalkPreorder(tree.Root, func(node *Node) {
        margin := ""
        p := node.Parent
        for p != tree.tnil {
            margin += "-"
            p = p.Parent
        }

        c := "Black"
        if node.Color == Red {
            c = "Red"
        }

        b.WriteString(fmt.Sprintf("\n%s %v (%d): %s", margin, getStringValueOf(node), node.Size, c))
    })

    // Assert
    t.Log(b.String())
}

func TestRbTreeSearchSuccess(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createTestRbTree()
    n := createStringNode("intel")

    // Act
    found, ok := Search(tree.Root, n)

    // Assert
    ass.True(ok)
    ass.NotNil(found)
    ass.Equal("intel", getStringValueOf(found))
}

func TestRbTreePredessor(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createTestRbTree()
    n := createStringNode("intel")
    found, _ := Search(tree.Root, n)

    // Act
    p := Predecessor(found)

    // Assert
    ass.Equal("fake", getStringValueOf(p))
}

func TestRbTreeSuccessor(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createTestRbTree()
    n := createStringNode("intel")
    found, _ := Search(tree.Root, n)

    // Act
    s := Successor(found)

    // Assert
    ass.Equal("it", getStringValueOf(s))
}

func TestRbTreeDelete(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createTestRbTree()
    n := createStringNode("intel")
    found, _ := Search(tree.Root, n)

    // Act
    DeleteFromRbTree(tree, found)

    // Assert
    found, ok := Search(tree.Root, n)
    ass.False(ok)
    ass.Nil(found.Key)

    found, ok = Search(tree.Root, createStringNode("microsoft"))
    ass.True(ok)
    ass.Equal("microsoft", getStringValueOf(found))
}

func createTestRbTree() *RbTree {
    nodes := []string{"abc", "amd", "cisco", "do", "fake", "intel", "it", "let", "microsoft", "russia", "usa", "xxx", "yyy", "zen"}
    tree := NewRbTree()
    for _, n := range nodes {
        RbTreeInsert(tree, &Node{Key: createStringNode(n)})
    }
    return tree
}

func createStringTestTree() *Node {
    r := createStringNode("fake")
    root := Node{Key: r}
    nodes := []string{"let", "zen", "yyy", "xxx", "do", "it", "amd", "intel", "cisco", "microsoft", "abc", "usa", "russia"}
    for _, n := range nodes {
        Insert(&root, &Node{Key: createStringNode(n)})
    }
    return &root
}

func getIntValueOf(node *Node) int {
    return (*node.Key).(IntNode).value
}

func getStringValueOf(node *Node) string {
    return (*node.Key).(StringNode).value
}

func createIntNode(v int) *Comparable {
    var r Comparable
    r = IntNode{v}
    return &r
}

func createStringNode(v string) *Comparable {
    var r Comparable
    r = StringNode{v}
    return &r
}
