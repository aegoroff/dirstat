package main

import (
	"github.com/aegoroff/dirstat/internal/cmd"
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/spf13/afero"
	"os"
)

func main() {
	if err := cmd.Execute(afero.NewOsFs(), out.NewConsoleEnvironment()); err != nil {
		os.Exit(1)
	}
}
