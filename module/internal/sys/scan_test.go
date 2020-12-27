package sys

import (
	"errors"
	"fmt"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

type testHandler struct {
	folders int
	files   int
	fipaths c9s.StringHashSet
	fopaths c9s.StringHashSet
	fp      []string
}

func (t *testHandler) Handle(evt *ScanEvent) {
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
	Scan("/", fs, &th)

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
	Scan("/", fs, &th)

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

type errCloser struct{}

func (e *errCloser) Close() error {
	return errors.New("new error")
}

func Test_Close_ThatReturnsError(t *testing.T) {
	// Arrange
	ec := &errCloser{}

	// Act
	Close(ec)

	// Assert
}
