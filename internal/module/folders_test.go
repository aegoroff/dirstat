package module

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/aegoroff/godatastruct/rbtree/special"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_bySizeFoldersTest(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	tree := special.NewMaxTree(2)

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

	tree.Insert(&fs1)
	tree.Insert(&fs2)
	tree.Insert(&fs3)

	// Assert
	var r []string
	rbtree.NewDescend(tree).Foreach(func(n rbtree.Comparable) {
		r = append(r, n.(folderI).String())
	})

	ass.ElementsMatch([]string{"/f3", "/f2"}, r)
}

func Test_byCountFoldersTest(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	tree := special.NewMaxTree(2)

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

	tree.Insert(&fc1)
	tree.Insert(&fc2)
	tree.Insert(&fc3)

	// Assert
	var r []string
	rbtree.NewDescend(tree).Foreach(func(n rbtree.Comparable) {
		r = append(r, n.(folderI).String())
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
	fc := newFolders(ctx.top)
	handler := newFoldersHandler(fc, ctx.pd)
	onlyFiles := newOnlyFoldersHandler(handler)

	// Act
	scan.Scan("/", newFs(appFS), onlyFiles)

	// Assert
	ass.Equal(int64(2), fc.byCount.Len())
	ass.Equal(int64(2), fc.bySize.Len())
}

func Test_folderHandler_EmptyFileSystem(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	ctx := NewContext(2, false, "/")
	fc := newFolders(ctx.top)
	handler := newFoldersHandler(fc, ctx.pd)
	onlyFiles := newOnlyFoldersHandler(handler)

	// Act
	scan.Scan("/", newFs(appFS), onlyFiles)

	// Assert
	ass.Equal(int64(1), fc.byCount.Len())
	ass.Equal(int64(1), fc.bySize.Len())
}

func Test_ExecuteFoldersModule_WithOutput(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/f/s", 0755)
	_ = afero.WriteFile(appFS, "/f/f.txt", []byte("123"), 0644)
	_ = afero.WriteFile(appFS, "/f/s/f.txt", []byte("1234"), 0644)
	ctx := NewContext(2, false, "/")
	m := NewFoldersModule(ctx, 0)
	e := out.NewMemoryEnvironment()
	p, _ := e.NewPrinter()

	// Act
	Execute("/", appFS, p, m)

	// Assert
	ass.Greater(len(e.String()), 0)
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
