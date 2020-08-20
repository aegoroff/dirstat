package main

import (
	"dirstat/cmd"
	"github.com/spf13/afero"
	"os"
)

func main() {
	if err := cmd.Execute(afero.NewOsFs(), os.Stdout); err != nil {
		os.Exit(1)
	}
}
