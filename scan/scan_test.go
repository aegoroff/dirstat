package scan

import (
	"errors"
	"fmt"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

type filesystem struct {
	fs afero.Fs
}

func newFs(fs afero.Fs) Filesystem {
	return &filesystem{fs: fs}
}

func (f *filesystem) Open(path string) (File, error) {
	return f.fs.Open(path)
}

type testHandler struct {
	folders int
	files   int
	fipaths c9s.StringHashSet
	fopaths c9s.StringHashSet
	fp      []string
}

func (t *testHandler) Handle(evt *Event) {
	if evt.File != nil {
		t.files++
		t.fipaths.Remove(evt.File.Path)
	}
	if evt.Folder != nil {
		t.folders++
		p := evt.Folder.Path
		t.fp = append(t.fp, p)
		t.fopaths.Remove(p)
	}
}

func Test_Scan(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/f/s", 0755)
	_ = afero.WriteFile(fs, "/f/f.txt", []byte("123"), 0644)
	_ = afero.WriteFile(fs, "/f/s/f.txt", []byte("1234"), 0644)

	th := testHandler{
		fipaths: make(c9s.StringHashSet),
		fopaths: make(c9s.StringHashSet),
		fp:      make([]string, 0),
	}

	// Act
	Scan("/", newFs(fs), &th)

	// Assert
	ass.Equal(2, th.files)
	ass.Equal(3, th.folders)
}

func Test_Scan_ManyData(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	fs := afero.NewMemMapFs()

	th := testHandler{
		fipaths: make(c9s.StringHashSet),
		fopaths: make(c9s.StringHashSet),
		fp:      make([]string, 0),
	}

	for i := 0; i < 1000; i++ {
		b := fmt.Sprintf("/f%d", i)
		c := filepath.Join(b, "s")
		th.fopaths.Add(b)
		th.fopaths.Add(c)

		_ = fs.MkdirAll(c, 0755)
		f1 := filepath.Join(b, "f.txt")
		f2 := filepath.Join(c, "f.txt")
		_ = afero.WriteFile(fs, f1, []byte("123"), 0644)
		_ = afero.WriteFile(fs, f2, []byte("1234"), 0644)
		th.fipaths.Add(f1)
		th.fipaths.Add(f2)
	}

	// Act
	Scan("/", newFs(fs), &th)

	// Assert
	ass.Equal(2000, th.files, "Invalid files count")
	ass.Equal(2001, th.folders, "Invalid folders count")

	fmt.Println("--- Missing files ---")
	for _, s := range th.fipaths.Items() {
		fmt.Println(s)
	}
	fmt.Println("--- Missing folders ---")
	for _, s := range th.fopaths.Items() {
		fmt.Println(s)
	}
	fm := make(map[string]int)
	for _, s := range th.fp {
		found, ok := fm[s]
		if ok {
			found++
			fm[s] = found
		} else {
			fm[s] = 1
		}
	}

	fmt.Println("--- Walked folders ---")
	for k, v := range fm {
		if v > 1 || k == "" {
			fmt.Printf("%s :%d\n", k, v)
		}
	}
}

func Test_Scan_OpenErrorsHandling(t *testing.T) {
	oefs := &testFs{
		err: errors.New("new error"),
	}

	tf := testFile{err: errors.New("new error")}
	refs := &testFs{
		f: &tf,
	}

	var tests = []struct {
		name string
		efs  Filesystem
	}{
		{"open error", oefs},
		{"readdir error", refs},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)

			th := testHandler{
				fipaths: make(c9s.StringHashSet),
				fopaths: make(c9s.StringHashSet),
				fp:      make([]string, 0),
			}

			// Act
			Scan("/", test.efs, &th)

			// Assert
			ass.Equal(0, th.files, "Invalid files count")
			ass.Equal(0, th.folders, "Invalid folders count")
		})
	}
}

type testFs struct {
	err error
	f   File
}

type testFile struct {
	err error
}

func (t *testFile) Close() error {
	return t.err
}

func (t *testFile) Readdir(_ int) ([]os.FileInfo, error) {
	return nil, t.err
}

func (e *testFs) Open(_ string) (File, error) {
	return e.f, e.err
}
