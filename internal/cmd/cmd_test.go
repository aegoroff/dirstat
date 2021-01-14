package cmd

import (
	"fmt"
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_PositiveTests(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	_ = appFS.MkdirAll("/f/s", 0755)
	_ = afero.WriteFile(appFS, "/f/f.txt", []byte("123"), 0644)
	_ = afero.WriteFile(appFS, "/f/s/f.md", []byte("1234"), 0644)
	_ = afero.WriteFile(appFS, "/f/f2.text", []byte("12345"), 0644)
	_ = afero.WriteFile(appFS, "/f/f3.xml", []byte("123456"), 0644)
	_ = afero.WriteFile(appFS, "/f/f31.xml", []byte("123456_"), 0644)
	_ = afero.WriteFile(appFS, "/f/f4.pdf", []byte("1234567"), 0644)
	_ = afero.WriteFile(appFS, "/f/f5.html", []byte("12345678"), 0644)
	_ = afero.WriteFile(appFS, "/f/f6.htm", []byte("12345678910"), 0644)

	var tests = []struct {
		name        string
		cmdline     []string
		mustcontain string
	}{
		{"a /", []string{"a", "/"}, "Total files"},
		{"a / -t 3", []string{"a", "/", "-t", "3"}, "TOP 3 file extensions by size"},
		{"a / -r 1", []string{"a", "/", "-r", "1"}, "1\\. Between 0 B and 100 KiB"},
		{"a / -r 1,2,3", []string{"a", "/", "-r", "1,2,3"}, "1\\. Between 0 B and 100 KiB"},
		{"a /f -r 1 -o", []string{"a", "/f", "-r", "1", "-o"}, "Total files"},
		{"a / -m", []string{"a", "/", "-m"}, "Total files"},
		{"a /f -o", []string{"a", "/f", "-o"}, "Total files"},
		{"b /", []string{"b", "/"}, "The first file's size digit distribution of non zero files \\(benford law\\)"},
		{"fi /", []string{"fi", "/"}, "TOP \\d+ files by size"},
		{"e /", []string{"e", "/"}, "TOP \\d+ file extensions by size"},
		{"fo /", []string{"fo", "/"}, "TOP \\d+ folders by size"},
		{"ver", []string{"ver"}, Version},
		{"nothing", []string{}, ""},
		{"fi :", []string{"fi", ":"}, ""},
		{"fi nothing", []string{"fi", ""}, ""},
		{"a / -p /res", []string{"a", "/", "-p", "/res"}, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := out.NewMemoryEnvironment()
			// Act
			err := Execute(appFS, w, test.cmdline...)

			// Assert
			o := w.String()
			ass.NoError(err)
			ass.Regexp(test.mustcontain, o)
			fmt.Println(strings.Join(test.cmdline, " "))
			fmt.Println("----------------------------------------------")
			fmt.Println(o)
		})
	}
}

func Test_NoCmdOptionSetTest(t *testing.T) {
	var tests = []struct {
		name    string
		cmdline []string
	}{
		{"a", []string{"a"}},
		{"b", []string{"b"}},
		{"fi", []string{"fi"}},
		{"fo", []string{"fo"}},
		{"e", []string{"e"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			appFS := afero.NewMemMapFs()
			w := out.NewMemoryEnvironment()

			// Act
			err := Execute(appFS, w, test.cmdline...)

			// Assert
			ass.NoError(err)
			ass.Equal("", w.String())
		})
	}
}

func Test_NegativeCmdTest(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	w := out.NewMemoryEnvironment()

	// Act
	err := Execute(appFS, w, "x")

	// Assert
	ass.Error(err)
}

func Test_ErrorToCreateOutputFile_Test(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	w := out.NewMemoryEnvironment()

	// Act
	err := Execute(afero.NewReadOnlyFs(appFS), w, "a", "/", "-p", "/out")

	// Assert
	ass.Error(err)
}

func Test_ConsoleEnvironment_Test(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	appFS := afero.NewMemMapFs()
	w := out.NewConsoleEnvironment()

	// Act
	err := Execute(appFS, w, "a", "/")

	// Assert
	ass.NoError(err)
}
