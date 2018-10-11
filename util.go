package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("\nAlloc = %s", humanize.IBytes(m.Alloc))
	fmt.Printf("\tTotalAlloc = %s", humanize.IBytes(m.TotalAlloc))
	fmt.Printf("\tSys = %s", humanize.IBytes(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func walkDirBreadthFirst(path string, action func(parent string, entry os.FileInfo)) {
	queue := make([]string, 0)

	queue = append(queue, path)

	for len(queue) > 0 {
		curr := queue[0]

		for _, entry := range dirents(curr) {
			action(curr, entry)
			if entry.IsDir() {
				queue = append(queue, filepath.Join(curr, entry.Name()))
			}
		}

		queue = queue[1:]
	}
}

func dirents(path string) []os.FileInfo {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	return entries
}
