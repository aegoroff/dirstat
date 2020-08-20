package main

import (
	"dirstat/cmd"
	"github.com/spf13/afero"
	"os"
)

func main() {
	cmd.Execute(afero.NewOsFs(), os.Stdout)
}
