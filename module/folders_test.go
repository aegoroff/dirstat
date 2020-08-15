package module

import (
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

	pd := pathDecorator{
		removeRoot: false,
		root:       "/",
	}
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
	tree.tree.Descend(func(n rbtree.Node) bool {
		r = append(r, n.Key().String())
		return true
	})

	ass.ElementsMatch([]string{"/f3", "/f2"}, r)
}

func Test_byCountFoldersTest(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	tree := newFixedTree(2)

	pd := pathDecorator{
		removeRoot: false,
		root:       "/",
	}
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
	tree.tree.Descend(func(n rbtree.Node) bool {
		r = append(r, n.Key().String())
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
	ass.Equal(int64(3), worker.total.CountFolders)
	ass.Equal(int64(2), worker.byCount.tree.Len())
	ass.Equal(int64(2), worker.bySize.tree.Len())
}
