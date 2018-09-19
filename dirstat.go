// This tool shows directory specified statistic.  This includes files and dirs count, size etc.
package main

import (
    "fmt"
    "github.com/dustin/go-humanize"
    "github.com/voxelbrain/goptions"
    "log"
    "os"
    "path/filepath"
    "sort"
    "text/tabwriter"
    "text/template"
    "time"
)

type options struct {
    Help      goptions.Help `goptions:"-h, --help, description='Show this help'"`
    Verbosity bool          `goptions:"-v, --verbose, description='Be verbose'"`
    Range     []int         `goptions:"-r, --range, description='Output verbose files info for fileSizeRanges specified'"`
    Path      string        `goptions:"-p, --path, obligatory, description='Name to the directory'"`
}

const (
    _ int64 = 1 << (10 * iota)
    Kbyte
    Mbyte
    Gbyte
    Tbyte
)

var fileSizeRanges = [...]Range{
    {Min: 0, Max: 100 * Kbyte},
    {Min: 100 * Kbyte, Max: Mbyte},
    {Min: Mbyte, Max: 10 * Mbyte},
    {Min: 10 * Mbyte, Max: 100 * Mbyte},
    {Min: 100 * Mbyte, Max: Gbyte},
    {Min: Gbyte, Max: 10 * Gbyte},
    {Min: 10 * Gbyte, Max: 100 * Gbyte},
    {Min: 100 * Gbyte, Max: Tbyte},
}

type fileStat struct {
    TotalFilesSize  uint64
    TotalFilesCount int64
}

type walkEntry struct {
    Size   int64
    Parent string
    Name   string
    IsDir  bool
}

type totalInfo struct {
    ReadingTime    time.Duration
    CountFiles     int64
    CountFolders   int64
    TotalFilesSize uint64
}

type namedInt64 struct {
    name  string;
    value int64
}

type namedInts64 []*namedInt64

func (x namedInts64) Len() int {
    return len(x)
}

func (x namedInts64) Less(i, j int) bool {
    return x[i].value < x[j].value
}

func (x namedInts64) Swap(i, j int) {
    x[i], x[j] = x[j], x[i]
}

func main() {
    opt := options{}

    goptions.ParseAndFail(&opt)

    if _, err := os.Stat(opt.Path); os.IsNotExist(err) {
        log.Fatalf("Directory '%s' does not exist. Details:\n  %v", opt.Path, err)
    }

    fmt.Printf("Root: %s\n\n", opt.Path)

    runAnalyze(opt)

    printMemUsage()
}

func runAnalyze(opt options) {
    total, stat, filesByRange, sizeByExt, countByExt := walk(opt)

    bySize := createSliceFromMap(sizeByExt)
    byCount := createSliceFromMap(countByExt)

    sort.Sort(sort.Reverse(bySize))
    sort.Sort(sort.Reverse(byCount))

    fmt.Print("Total files stat:\n\n")

    const format = "%v\t%v\t%v\t%v\t%v\n"
    tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

    fmt.Fprintf(tw, format, "File size", "Amount", "%", "Size", "%")
    fmt.Fprintf(tw, format, "---------", "------", "------", "----", "------")

    var heads []string
    for i, r := range fileSizeRanges {
        percentOfCount := (float64(stat[r].TotalFilesCount) / float64(total.CountFiles)) * 100
        percentOfSize := (float64(stat[r].TotalFilesSize) / float64(total.TotalFilesSize)) * 100
        head := fmt.Sprintf("%d. Between %s and %s", i+1, humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)))
        heads = append(heads, head)

        fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", head, stat[r].TotalFilesCount, percentOfCount, humanize.IBytes(stat[r].TotalFilesSize), percentOfSize)
    }
    tw.Flush()

    fmt.Print("\nTOP 10 file extensions by size:\n\n")
    fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
    fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

    for i := 0; i < 10; i++ {
        h := bySize[i].name
        percentOfCount := (float64(countByExt[h]) / float64(total.CountFiles)) * 100
        percentOfSize := (float64(bySize[i].value) / float64(total.TotalFilesSize)) * 100

        sz := uint64(bySize[i].value)
        fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", h, countByExt[h], percentOfCount, humanize.IBytes(sz), percentOfSize)
    }

    tw.Flush()

    fmt.Print("\nTOP 10 file extensions by count:\n\n")
    fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
    fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

    for i := 0; i < 10; i++ {
        h := byCount[i].name
        percentOfCount := (float64(byCount[i].value) / float64(total.CountFiles)) * 100
        percentOfSize := (float64(sizeByExt[h]) / float64(total.TotalFilesSize)) * 100

        sz := uint64(sizeByExt[h])
        fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", h, byCount[i].value, percentOfCount, humanize.IBytes(sz), percentOfSize)
    }

    tw.Flush()

    if opt.Verbosity && len(opt.Range) > 0 {
        fmt.Printf("\nDetailed files stat:\n")
        for i, r := range fileSizeRanges {
            if len(filesByRange[r]) == 0 {
                continue
            }

            fmt.Printf("%s\n", heads[i])
            for _, item := range filesByRange[r] {
                fullPath := filepath.Join(item.Parent, item.Name)
                size := humanize.IBytes(uint64(item.Size))
                fmt.Printf("   %s - %s\n", fullPath, size)
            }
        }
    }

    printTotals(total)
}

func createSliceFromMap(sizeByExt map[string]int64) namedInts64 {
    var result = make(namedInts64, len(sizeByExt))
    i := 0
    for k, v := range sizeByExt {
        result[i] = &namedInt64{value: v, name: k}
        i++
    }
    return result
}

func walk(opt options) (totalInfo, map[Range]fileStat, map[Range][]*walkEntry, map[string]int64, map[string]int64) {
    verboseRanges := make(map[int]bool)
    for _, x := range opt.Range {
        verboseRanges[x] = true
    }
    total := totalInfo{}
    stat := make(map[Range]fileStat)
    sizeByExt := make(map[string]int64)
    countByExt := make(map[string]int64)
    filesByRange := map[Range][]*walkEntry{}

    ch := make(chan *walkEntry, 1024)

    start := time.Now()

    go func(ch chan<- *walkEntry) {
        walkDirBreadthFirst(opt.Path, func(parent string, entry os.FileInfo) {
            ch <- &walkEntry{IsDir: entry.IsDir(), Size: entry.Size(), Parent: parent, Name: entry.Name()}
        })
        close(ch)
    }(ch)

    for {
        we, ok := <-ch
        if !ok {
            break
        }

        if we.IsDir {
            total.CountFolders++
        } else {
            // Accumulate file statistic
            sz := uint64(we.Size)
            total.CountFiles++
            total.TotalFilesSize += sz
            for i, r := range fileSizeRanges {
                if !r.contains(float64(sz)) {
                    continue
                }

                s := stat[r]
                s.TotalFilesCount++
                s.TotalFilesSize += sz
                stat[r] = s

                ext := filepath.Ext(we.Name)
                sizeByExt[ext] += we.Size
                countByExt[ext]++

                if !opt.Verbosity || !verboseRanges[i+1] {
                    continue
                }

                nodes, ok := filesByRange[r]
                if !ok {
                    filesByRange[r] = []*walkEntry{we}
                } else {
                    filesByRange[r] = append(nodes, we)
                }
            }
        }
    }

    total.ReadingTime = time.Since(start)
    return total, stat, filesByRange, sizeByExt, countByExt
}

func printTotals(t totalInfo) {

    const totalTemplate = `
Total files:   {{.CountFiles}} ({{.TotalFilesSize | toBytesString }})
Total folders: {{.CountFolders}}

Read taken:    {{.ReadingTime}}
`

    var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
    report.Execute(os.Stdout, t)
}
