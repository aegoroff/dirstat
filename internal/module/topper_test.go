package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/aegoroff/godatastruct/rbtree/special"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_topperPrintInvalidTree_invalidCastingError(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	ft := special.NewMaxTree(3)
	ft.Insert(rbtree.Int(1))
	ft.Insert(rbtree.Int(2))
	ft.Insert(rbtree.Int(3))

	e := out.NewMemoryEnvironment()
	p, _ := e.NewPrinter()
	tp := newTopper(p, nil, []string{"#", "Extension", "Count", "%", "Size", "%"}, &noChangeDecorator{})

	// Act
	tp.descend(ft)

	// Assert
	ass.Contains(e.String(), "Invalid casting: expected folderI key type but it wasn`t")
}
