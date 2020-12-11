package module

import (
	"bytes"
	"dirstat/module/internal/sys"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_bySizeFoldersTest(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	tree := newFixedTree(2)

	pd := nonDestructiveDecorator{}
	f1 := folder{
		path:  "/f1",
		count: 3,
		size:  100,
		pd:    &pd,
	}
	fs1 := folderS{f1}
	f2 := folder{
		path:  "/f2",
		count: 2,
		size:  200,
		pd:    &pd,
	}
	fs2 := folderS{f2}
	f3 := folder{
		path:  "/f3",
		count: 1,
		size:  300,
		pd:    &pd,
	}
	fs3 := folderS{f3}

	tree.insert(&fs1)
	tree.insert(&fs2)
	tree.insert(&fs3)

	// Assert
	var r []string
	tree.descend(func(n rbtree.Node) bool {
		r = append(r, n.String())
		return true
	})

	ass.ElementsMatch([]string{"/f3", "/f2"}, r)
}

func Test_byCountFoldersTest(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	tree := newFixedTree(2)

	pd := nonDestructiveDecorator{}
	f1 := folder{
		path:  "/f1",
		count: 3,
		size:  100,
		pd:    &pd,
	}
	fc1 := folderC{f1}
	f2 := folder{
		path:  "/f2",
		count: 2,
		size:  200,
		pd:    &pd,
	}
	fc2 := folderC{f2}
	f3 := folder{
		path:  "/f3",
		count: 1,
		size:  300,
		pd:    &pd,
	}
	fc3 := folderC{f3}

	tree.insert(&fc1)
	tree.insert(&fc2)
	tree.insert(&fc3)

	// Assert
	var r []string
	tree.descend(func(n rbtree.Node) bool {
		r = append(r, n.String())
		return true
	})

	ass.ElementsMatch([]string{"/f1", "/f2"}, r)
}

func Test_folderHandler(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/f/s", 0755)
	_ = afero.WriteFile(appFS, "/f/f.txt", []byte("123"), 0644)
	_ = afero.WriteFile(appFS, "/f/s/f.txt", []byte("1234"), 0644)
	ctx := NewContext(2, false, "/")
	worker := newFoldersWorker(ctx)
	var handlers []sys.ScanHandler
	handlers = append(handlers, worker.handler)

	// Act
	sys.Scan("/", appFS, handlers)

	// Assert
	ass.Equal(int64(2), worker.byCount.tree.Len())
	ass.Equal(int64(2), worker.bySize.tree.Len())
}

func Test_folderHandler_EmptyFileSystem(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	ctx := NewContext(2, false, "/")
	worker := newFoldersWorker(ctx)
	var handlers []sys.ScanHandler
	handlers = append(handlers, worker.handler)

	// Act
	sys.Scan("/", appFS, handlers)

	// Assert
	ass.Equal(int64(1), worker.byCount.tree.Len())
	ass.Equal(int64(1), worker.bySize.tree.Len())
}

func Test_ExecuteFoldersModule_WithOutput(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/f/s", 0755)
	_ = afero.WriteFile(appFS, "/f/f.txt", []byte("123"), 0644)
	_ = afero.WriteFile(appFS, "/f/s/f.txt", []byte("1234"), 0644)
	ctx := NewContext(2, false, "/")
	m := NewFoldersModule(ctx)
	w := bytes.NewBufferString("")

	// Act
	Execute("/", appFS, w, m)

	// Assert
	ass.Greater(w.Len(), 0)
}

func TestFolder_Path(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	fo := folder{
		path:  "/usr",
		size:  0,
		count: 0,
		pd:    &nonDestructiveDecorator{},
	}

	// Act
	result := fo.Path()

	// Assert
	ass.Equal("/usr", result)
}

func TestFolder_Path_PathDecorating(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	pd := &removeRootDecorator{
		root: "/usr",
	}
	fo := folder{
		path:  "/usr/local",
		size:  0,
		count: 0,
		pd:    pd,
	}

	// Act
	result := fo.Path()

	// Assert
	ass.Equal("/local", result)
}

func Test_castSize_invalidCasting(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	f := folder{
		path:  "/f1",
		count: 3,
		size:  100,
		pd:    nil,
	}
	fc := folderC{f}

	// Act
	r, err := castSize(&fc)

	// Assert
	ass.Error(err)
	ass.Nil(r)
}

func Test_castCount_invalidCasting(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	f := folder{
		path:  "/f1",
		count: 3,
		size:  100,
		pd:    nil,
	}
	fs := folderS{f}

	// Act
	r, err := castCount(&fs)

	// Assert
	ass.Error(err)
	ass.Nil(r)
}

func Test_printTop_invalidCastingError(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	ft := newFixedTree(3)
	ft.insert(rbtree.NewInt(1))
	ft.insert(rbtree.NewInt(2))
	ft.insert(rbtree.NewInt(3))
	fr := foldersRenderer{}
	w := bytes.NewBufferString("")

	// Act
	fr.printTop(ft, newPrinter(w), castSize)

	// Assert
	ass.Contains(w.String(), "invalid casting: expected *folderS key type but it wasn`t")
}
