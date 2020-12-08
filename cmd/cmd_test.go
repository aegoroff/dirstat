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
		name        string
		cmdline     []string
		mustcontain string
	}{
		{"a -p /", []string{"a", "-p", "/"}, "Total files"},
		{"a -p / -t 3", []string{"a", "-p", "/", "-t", "3"}, "Total files"},
		{"a -p / -r 1", []string{"a", "-p", "/", "-r", "1"}, "Total files"},
		{"a -p /f -r 1 -o", []string{"a", "-p", "/f", "-r", "1", "-o"}, "Total files"},
		{"a -p / -m", []string{"a", "-p", "/", "-m"}, "Total files"},
		{"a -p /f -o", []string{"a", "-p", "/f", "-o"}, "Total files"},
		{"b -p /", []string{"b", "-p", "/"}, "Total files"},
		{"fi -p /", []string{"fi", "-p", "/"}, "Total files"},
		{"e -p /", []string{"e", "-p", "/"}, "Total files"},
		{"fo -p /", []string{"fo", "-p", "/"}, "Total files"},
		{"ver", []string{"ver"}, "dirstat v"},
		{"nothing", []string{}, ""},
		{"fi -p :", []string{"fi", "-p", ":"}, ""},
		{"fi -p nothing", []string{"fi", "-p", ""}, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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
		})
	}
}

func Test_NegativeCmdTest(t *testing.T) {
	var tests = []struct {
		name    string
		cmdline []string
	}{
		{"a", []string{"a"}},
		{"b", []string{"b"}},
		{"fi", []string{"fi"}},
		{"fo", []string{"fo"}},
		{"e", []string{"e"}},
		{"x", []string{"x"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			appFS := afero.NewMemMapFs()
			w := bytes.NewBufferString("")

			// Act
			err := Execute(appFS, w, test.cmdline...)

			// Assert
			ass.Error(err)
			fmt.Println(w.String())
		})
	}
}
