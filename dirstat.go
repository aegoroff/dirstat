// This tool shows directory specified statistic.  This includes files and dirs count, size etc.
package main

import (
    "dirstat/tree"
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

const Top = 10

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
    ReadingTime   time.Duration
    FilesTotal    countSizeAggregate
    CountFolders  int64
    CountFileExts int
}

type countSizeAggregate struct {
    Count int64
    Size  uint64
}

type namedInt64 struct {
    name  string
    value int64
}

type statItem struct {
    name  string
    size  int64
    count int64
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

func (x statItem) LessThan(y interface{}) bool {
    return x.size < (y.(statItem)).size
}

func (x statItem) EqualTo(y interface{}) bool {
    return x.size == (y.(statItem)).size
}

func (x namedInt64) LessThan(y interface{}) bool {
    return x.value < (y.(namedInt64)).value
}

func (x namedInt64) EqualTo(y interface{}) bool {
    return x.value == (y.(namedInt64)).value
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
    total, stat, filesByRange, byExt, byFolder, topFiles := walk(opt)
    total.CountFileExts = len(byExt)

    extBySize := createSliceFromMap(byExt, func(aggregate countSizeAggregate) int64 {
        return int64(aggregate.Size)
    })

    extByCount := createSliceFromMap(byExt, func(aggregate countSizeAggregate) int64 {
        return aggregate.Count
    })

    sort.Sort(sort.Reverse(extBySize))
    sort.Sort(sort.Reverse(extByCount))

    fmt.Print("Total files stat:\n\n")

    const format = "%v\t%v\t%v\t%v\t%v\n"
    tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

    fmt.Fprintf(tw, format, "File size", "Amount", "%", "Size", "%")
    fmt.Fprintf(tw, format, "---------", "------", "------", "----", "------")

    var heads []string
    for i, r := range fileSizeRanges {
        h := fmt.Sprintf("%d. Between %s and %s", i+1, humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)))
        heads = append(heads, h)

        count := stat[r].TotalFilesCount
        sz := stat[r].TotalFilesSize

        percentOfCount := countPercent(count, total)
        percentOfSize := sizePercent(sz, total)

        fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", h, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
    }
    tw.Flush()

    fmt.Printf("\nTOP %d file extensions by size:\n\n", Top)
    fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
    fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

    for i := 0; i < Top; i++ {
        h := extBySize[i].name

        count := byExt[h].Count
        sz := uint64(extBySize[i].value)

        percentOfCount := countPercent(count, total)
        percentOfSize := sizePercent(sz, total)

        fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", h, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
    }

    tw.Flush()

    fmt.Printf("\nTOP %d file extensions by count:\n\n", Top)
    fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
    fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

    for i := 0; i < Top; i++ {
        h := extByCount[i].name

        count := extByCount[i].value
        sz := byExt[h].Size

        percentOfCount := countPercent(count, total)
        percentOfSize := sizePercent(sz, total)

        fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", h, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
    }

    tw.Flush()

    fmt.Printf("\nTOP %d files by size:\n\n", Top)
    fmt.Fprintf(tw, "%v\t%v\n", "File", "Size")
    fmt.Fprintf(tw, "%v\t%v\n", "------", "----")

    fileTreeSize := topFiles.Root.Size
    for i := fileTreeSize; i > 0; i-- {
        n := tree.OrderStatisticSelect(topFiles.Root, i)
        order := (*n.Key).(namedInt64)
        h := fmt.Sprintf("%d. %s", fileTreeSize-i+1, order.name)

        sz := uint64(order.value)

        fmt.Fprintf(tw, "%v\t%v\n", h, humanize.IBytes(sz))
    }

    tw.Flush()

    fmt.Printf("\nTOP %d folders by size:\n\n", Top)
    fmt.Fprintf(tw, format, "Folder", "Files", "%", "Size", "%")
    fmt.Fprintf(tw, format, "------", "-----", "------", "----", "------")

    treeSize := byFolder.Root.Size
    for i := treeSize; i > 0; i-- {
        n := tree.OrderStatisticSelect(byFolder.Root, i)
        order := (*n.Key).(statItem)
        h := fmt.Sprintf("%d. %s", treeSize-i+1, order.name)

        count := order.count
        sz := uint64(order.size)

        percentOfCount := countPercent(count, total)
        percentOfSize := sizePercent(sz, total)

        fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", h, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
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

func countPercent(count int64, total totalInfo) float64 {
    return (float64(count) / float64(total.FilesTotal.Count)) * 100
}

func sizePercent(size uint64, total totalInfo) float64 {
    return (float64(size) / float64(total.FilesTotal.Size)) * 100
}

func createSliceFromMap(sizeByExt map[string]countSizeAggregate, mapper func(countSizeAggregate) int64) namedInts64 {
    var result = make(namedInts64, len(sizeByExt))
    i := 0
    for k, v := range sizeByExt {
        result[i] = &namedInt64{value: mapper(v), name: k}
        i++
    }
    return result
}

func walk(opt options) (totalInfo, map[Range]fileStat, map[Range][]*walkEntry, map[string]countSizeAggregate, *tree.RbTree, *tree.RbTree) {
    verboseRanges := make(map[int]bool)
    for _, x := range opt.Range {
        verboseRanges[x] = true
    }
    total := totalInfo{}
    stat := make(map[Range]fileStat)
    filesByRange := make(map[Range][]*walkEntry)

    byExt := make(map[string]countSizeAggregate)

    ch := make(chan *walkEntry, 1024)

    start := time.Now()

    go func(ch chan<- *walkEntry) {
        walkDirBreadthFirst(opt.Path, func(parent string, entry os.FileInfo) {
            ch <- &walkEntry{IsDir: entry.IsDir(), Size: entry.Size(), Parent: parent, Name: entry.Name()}
        })
        close(ch)
    }(ch)

    folderSizeTree := tree.NewRbTree()
    fileSizeTree := tree.NewRbTree()

    currFolderStat := statItem{}

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
            total.FilesTotal.Count++
            total.FilesTotal.Size += sz

            ext := filepath.Ext(we.Name)
            a := byExt[ext]
            a.Size += uint64(we.Size)
            a.Count++
            byExt[ext] = a

            if currFolderStat.name == "" {
                currFolderStat.name = we.Parent
            }

            if currFolderStat.name == we.Parent {
                currFolderStat.size += we.Size
                currFolderStat.count++
            } else {
                if folderSizeTree.Root == nil || folderSizeTree.Root.Size < Top {
                    node := createTreeNode(currFolderStat)
                    tree.Insert(folderSizeTree, node)
                } else {
                    minSizeNode := tree.Minimum(folderSizeTree.Root)
                    if getSizeFromNode(minSizeNode) < currFolderStat.size {
                        tree.Delete(folderSizeTree, minSizeNode)

                        node := createTreeNode(currFolderStat)
                        tree.Insert(folderSizeTree, node)
                    }
                }
                currFolderStat.name = we.Parent
                currFolderStat.count = 1
                currFolderStat.size = we.Size
            }

            if fileSizeTree.Root == nil || fileSizeTree.Root.Size < Top {
                fullPath := filepath.Join(we.Parent, we.Name)
                node := createTreeNode(namedInt64{value: we.Size, name: fullPath})
                tree.Insert(fileSizeTree, node)
            } else {
                minSizeNode := tree.Minimum(fileSizeTree.Root)
                if getSizeFromNode(minSizeNode) < we.Size {
                    tree.Delete(fileSizeTree, minSizeNode)

                    fullPath := filepath.Join(we.Parent, we.Name)
                    node := createTreeNode(namedInt64{value: we.Size, name: fullPath})
                    tree.Insert(fileSizeTree, node)
                }
            }

            for i, r := range fileSizeRanges {
                if !r.contains(float64(sz)) {
                    continue
                }

                s := stat[r]
                s.TotalFilesCount++
                s.TotalFilesSize += sz
                stat[r] = s

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
    return total, stat, filesByRange, byExt, folderSizeTree, fileSizeTree
}

func createTreeNode(si tree.Comparable) *tree.Node {
    return &tree.Node{Key: &si}
}

func getSizeFromNode(node *tree.Node) int64 {
    if k, ok := (*node.Key).(statItem); ok {
        return k.size
    }

    if k, ok := (*node.Key).(namedInt64); ok {
        return k.value
    }

    return 0
}

func printTotals(t totalInfo) {

    const totalTemplate = `
Total files:            {{.FilesTotal.Count}} ({{.FilesTotal.Size | toBytesString }})
Total folders:          {{.CountFolders}}
Total file extensions:  {{.CountFileExts}}

Read taken:    {{.ReadingTime}}
`

    var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
    report.Execute(os.Stdout, t)
}
