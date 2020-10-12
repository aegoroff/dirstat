package cmd

import (
	"bytes"
	"fmt"
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
	_ = afero.WriteFile(appFS, "/f/f4.pdf", []byte("1234567"), 0644)
	_ = afero.WriteFile(appFS, "/f/f5.html", []byte("12345678"), 0644)
	_ = afero.WriteFile(appFS, "/f/f6.htm", []byte("12345678910"), 0644)

	var tests = []struct {
		cmdline     []string
		mustcontain string
	}{
		{[]string{"a", "-p", "/"}, "Total files"},
		{[]string{"a", "-p", "/", "-t", "3"}, "Total files"},
		{[]string{"a", "-p", "/", "-r", "1"}, "Total files"},
		{[]string{"a", "-p", "/f", "-r", "1", "-o"}, "Total files"},
		{[]string{"a", "-p", "/", "-m"}, "Total files"},
		{[]string{"a", "-p", "/f", "-o"}, "Total files"},
		{[]string{"b", "-p", "/"}, "Total files"},
		{[]string{"fi", "-p", "/"}, "Total files"},
		{[]string{"e", "-p", "/"}, "Total files"},
		{[]string{"fo", "-p", "/"}, "Total files"},
		{[]string{"ver"}, "dirstat v"},
		{[]string{}, ""},
		{[]string{"fi", "-p", ":"}, ""},
	}

	for _, test := range tests {
		w := bytes.NewBufferString("")
		// Act
		err := Execute(appFS, w, test.cmdline...)

		// Assert
		out := w.String()
		ass.NoError(err)
		ass.Contains(out, test.mustcontain)
		fmt.Println(strings.Join(test.cmdline, " "))
		fmt.Println("----------------------------------------------")
		fmt.Println(out)
	}
}

func Test_NegativeCmdTest(t *testing.T) {
	var tests = []struct {
		cmdline []string
	}{
		{[]string{"a"}},
		{[]string{"b"}},
		{[]string{"fi"}},
		{[]string{"fo"}},
		{[]string{"e"}},
		{[]string{"x"}},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		appFS := afero.NewMemMapFs()
		w := bytes.NewBufferString("")

		// Act
		err := Execute(appFS, w, test.cmdline...)

		// Assert
		ass.Error(err)
		fmt.Println(w.String())
	}
}
