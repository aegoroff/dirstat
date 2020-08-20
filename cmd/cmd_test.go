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
	_ = afero.WriteFile(appFS, "/f/s/f.txt", []byte("1234"), 0644)
	_ = afero.WriteFile(appFS, "/f/f2.txt", []byte("12345"), 0644)
	_ = afero.WriteFile(appFS, "/f/f3.txt", []byte("123456"), 0644)
	_ = afero.WriteFile(appFS, "/f/f4.txt", []byte("1234567"), 0644)
	_ = afero.WriteFile(appFS, "/f/f5.txt", []byte("12345678"), 0644)
	_ = afero.WriteFile(appFS, "/f/f6.txt", []byte("12345678910"), 0644)

	var tests = []struct {
		cmdline []string
	}{
		{[]string{"a", "-p", "/"}},
		{[]string{"a", "-p", "/", "-r", "1"}},
		{[]string{"a", "-p", "/f", "-r", "1", "-o"}},
		{[]string{"a", "-p", "/", "-m"}},
		{[]string{"a", "-p", "/f", "-o"}},
		{[]string{"b", "-p", "/"}},
		{[]string{"fi", "-p", "/"}},
		{[]string{"fi", "-p", "/", "-e"}},
		{[]string{"fo", "-p", "/"}},
		{[]string{"ver"}},
		{[]string{}},
	}

	for _, test := range tests {
		w := bytes.NewBufferString("")
		// Act
		err := Execute(appFS, w, test.cmdline...)

		// Assert
		ass.NoError(err)
		fmt.Println(strings.Join(test.cmdline, " "))
		fmt.Println("----------------------------------------------")
		fmt.Println(w.String())
	}
}
