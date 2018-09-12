package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
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
