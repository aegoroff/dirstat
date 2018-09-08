// This tool shows directory specified statistic.  This includes files and dirs count, size etc.
package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/voxelbrain/goptions"
	"log"
	"os"
	"path/filepath"
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

	countFiles := 0
	countDirs := 0
	var totalSize uint64
	filepath.Walk(options.Path, func(path string, info os.FileInfo, err error) error {
		fmt.Printf("%s", path)
		if info.IsDir() {
			countDirs++
		} else {
			countFiles++
			sz := uint64(info.Size())
			fmt.Printf(" (%s)", humanize.Bytes(sz))
			totalSize += sz
		}
		fmt.Print("\n")
		return nil
	})

	fmt.Printf("Total files %d Total size: %s\n", countFiles, humanize.Bytes(totalSize))
	fmt.Printf("Total folders %d\n", countDirs)
}
