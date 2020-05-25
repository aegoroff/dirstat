package cmd

import (
	"time"
)

const (
	_ int64 = 1 << (10 * iota)
	Kbyte
	Mbyte
	Gbyte
	Tbyte
)

const Top = 10

type options struct {
	Verbosity bool
	Range     []int
	Path      string
}

var fileSizeRanges = [...]Range{
	{Min: 0, Max: 100 * Kbyte},
	{Min: 100 * Kbyte, Max: Mbyte},
	{Min: Mbyte, Max: 10 * Mbyte},
	{Min: 10 * Mbyte, Max: 100 * Mbyte},
	{Min: 100 * Mbyte, Max: Gbyte},
	{Min: Gbyte, Max: 10 * Gbyte},
	{Min: 10 * Gbyte, Max: 100 * Gbyte},
	{Min: 100 * Gbyte, Max: Tbyte},
	{Min: Tbyte, Max: 10 * Tbyte},
	{Min: 10 * Tbyte, Max: 100 * Tbyte},
}

type fileStat struct {
	TotalFilesSize  uint64
	TotalFilesCount int64
}

type fileEntry struct {
	Size   int64
	Parent string
	Name   string
}

type totalInfo struct {
	ReadingTime   time.Duration
	FilesTotal    countSizeAggregate
	CountFolders  int64
	CountFileExts int
}

type countSizeAggregate struct {
	Count int64
	Size  uint64
}
