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

func Test_InorderWalkInt_AllElementsAscending(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()
    var result []int

    // Act
    WalkInorder(tree.Root, func(node *Node) {
        result = append(result, getIntValueOf(node))
    })

    // Assert
    ass.ElementsMatch(result, []int{2, 3, 4, 6, 7, 9, 13, 15, 17, 18, 20})
}

func Test_InorderWalkString_AllElementsAscending(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createTestStringTree()
    var result []string

    // Act
    WalkInorder(tree.Root, func(node *Node) {
        result = append(result, getStringValueOf(node))
    })

    // Assert
    ass.ElementsMatch(result, []string{"abc", "amd", "cisco", "do", "fake", "intel", "it", "let", "microsoft", "russia", "usa", "xxx", "yyy", "zen"})
}

func Test_OrderStatisticSelect_ValueAsExpected(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()

    // Act
    found := OrderStatisticSelect(tree.Root, 2)

    // Assert
    ass.NotNil(found)
    ass.Equal(3, getIntValueOf(found))
}

func Test_SearchIntTree_Success(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()
    v := createIntNode(13)

    // Act
    found, ok := Search(tree.Root, v)

    // Assert
    ass.True(ok)
    ass.NotNil(found)
    ass.Equal(13, getIntValueOf(found))
}

func Test_SearchStringTree_Success(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createTestStringTree()
    n := createStringNode("intel")

    // Act
    found, ok := Search(tree.Root, n)

    // Assert
    ass.True(ok)
    ass.NotNil(found)
    ass.Equal("intel", getStringValueOf(found))
}

func Test_SearchIntTree_Failure(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()
    v := createIntNode(22)

    // Act
    found, ok := Search(tree.Root, v)

    // Assert
    ass.False(ok)
    ass.Nil(found)
}

func Test_SuccessorInTheMiddle_ReturnSuccessor(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()
    v := createIntNode(13)
    r, _ := Search(tree.Root, v)

    // Act
    s := Successor(r)

    // Assert
    ass.Equal(15, getIntValueOf(s))
}

func Test_SuccessorOfMax_ReturnNil(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()
    v := createIntNode(20)
    r, _ := Search(tree.Root, v)

    // Act
    s := Successor(r)

    // Assert
    ass.Nil(s)
}

func Test_PredecessorInTheMiddle_PredecessorFound(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()
    v := createIntNode(13)
    r, _ := Search(tree.Root, v)

    // Act
    p := Predecessor(r)

    // Assert
    ass.Equal(9, getIntValueOf(p))
}

func Test_PredecessorOfMin_ReturnNil(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()
    v := createIntNode(2)
    r, _ := Search(tree.Root, v)

    // Act
    p := Predecessor(r)

    // Assert
    ass.Nil(p)
}

func Test_Minimum_ValueAsExpected(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()

    // Act
    r := Minimum(tree.Root)

    // Assert
    ass.Equal(2, getIntValueOf(r))
}

func Test_Maximum_ValueAsExpected(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createIntegerTestTree()

    // Act
    r := Maximum(tree.Root)

    // Assert
    ass.Equal(20, getIntValueOf(r))
}

func Test_RightRotate_StructureAsExpected(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    r := Node{Key: createStringNode("root")}

    tree := NewRbTree()
    Insert(tree, &r)

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

func Test_LeftRotate_StructureAsExpected(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    r := Node{Key: createStringNode("root")}

    tree := NewRbTree()
    Insert(tree, &r)

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

func Test_WalkPreorder(t *testing.T) {
    // Arrange
    tree := createTestStringTree()
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

func Test_Delete_NodeDeleted(t *testing.T) {
    // Arrange
    ass := assert.New(t)
    tree := createTestStringTree()
    n := createStringNode("intel")
    found, _ := Search(tree.Root, n)

    // Act
    Delete(tree, found)

    // Assert
    found, ok := Search(tree.Root, n)
    ass.False(ok)
    ass.Nil(found)

    found, ok = Search(tree.Root, createStringNode("microsoft"))
    ass.True(ok)
    ass.Equal("microsoft", getStringValueOf(found))
}

func createIntegerTestTree() *RbTree {
    nodes := []int{6, 18, 3, 15, 7, 2, 4, 13, 9, 17, 20}
    tree := NewRbTree()
    for _, n := range nodes {
        Insert(tree, &Node{Key: createIntNode(n)})
    }
    return tree
}

func createTestStringTree() *RbTree {
    nodes := []string{"abc", "amd", "cisco", "do", "fake", "intel", "it", "let", "microsoft", "russia", "usa", "xxx", "yyy", "zen"}
    tree := NewRbTree()
    for _, n := range nodes {
        Insert(tree, &Node{Key: createStringNode(n)})
    }
    return tree
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
