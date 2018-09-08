// This tool shows directory specified statistic.  This includes files and dirs count, size etc.
package main

import (
	"github.com/voxelbrain/goptions"
	"log"
	"os"
)

type Options struct {
	Help goptions.Help `goptions:"-h, --help, description='Show this help'"`
	Path string        `goptions:"-p, --path, obligatory, description='Path to the directory'"`
}

func main() {
	options := Options{}

	goptions.ParseAndFail(&options)

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		log.Fatalf("Directory '%s' does not exist. Details:\n  %v", options.Path, err)
	}
}
