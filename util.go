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

func walkDir(path string, action func(path string, info os.FileInfo)) {
	for _, entry := range dirents(path) {
		fullPath := filepath.Join(path, entry.Name())
		action(fullPath, entry)
		if entry.IsDir() {
			walkDir(fullPath, action)
		}
	}
}

func dirents(path string) []os.FileInfo {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	return entries
}
