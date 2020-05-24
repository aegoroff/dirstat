package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"text/tabwriter"
	"text/template"
	"time"
)

const pathParamName = "path"
const verboseParamName = "verbose"
const rangeParamName = "range"

var appFileSystem = afero.NewOsFs()
var appWriter io.Writer

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "dirstat",
	Short: "Directory statistic tool",
	Long:  ` A small tool that shows selected folder or drive (on Windows) usage statistic`,
	RunE: func(cmd *cobra.Command, args []string) error {

		path, err := cmd.Flags().GetString(pathParamName)

		if err != nil {
			return err
		}

		ranges, err := cmd.Flags().GetIntSlice(rangeParamName)

		if err != nil {
			return err
		}

		verbose, err := cmd.Flags().GetBool(verboseParamName)

		if err != nil {
			return err
		}

		opt := options{Verbosity: verbose, Path: path, Range: ranges}

		if _, err := appFileSystem.Stat(opt.Path); os.IsNotExist(err) {
			return err
		}

		if opt.Path[len(opt.Path)-1] == ':' {
			opt.Path = filepath.Join(opt.Path, "\\")
		}

		_, _ = fmt.Fprintf(appWriter, "Root: %s\n\n", opt.Path)

		runAnalyze(opt, appFileSystem, appWriter)

		printMemUsage(appWriter)
		return nil
	},
}

func init() {
	cobra.MousetrapHelpText = ""
	appWriter = os.Stdout
	rootCmd.Flags().StringP(pathParamName, "p", "", "REQUIRED. Directory path to show info.")
	rootCmd.Flags().IntSliceP(rangeParamName, "r", []int{}, "Output verbose files info for fileSizeRanges specified")
	rootCmd.Flags().BoolP(verboseParamName, "v", false, "Be verbose")
}

// Execute starts package running
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runAnalyze(opt options, fs afero.Fs, w io.Writer) {
	total, stat, filesByRange, byExt, topFolders, topFiles := walk(opt, fs)
	total.CountFileExts = len(byExt)

	extBySize := createSliceFromMap(byExt, func(aggregate countSizeAggregate) int64 {
		return int64(aggregate.Size)
	})

	extByCount := createSliceFromMap(byExt, func(aggregate countSizeAggregate) int64 {
		return aggregate.Count
	})

	sort.Sort(sort.Reverse(extBySize))
	sort.Sort(sort.Reverse(extByCount))

	_, _ = fmt.Fprintf(w, "Total files stat:\n\n")

	const format = "%v\t%v\t%v\t%v\t%v\n"
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	_, _ = fmt.Fprintf(tw, format, "File size", "Amount", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "------", "------", "----", "------")

	var heads []string
	for i, r := range fileSizeRanges {
		h := fmt.Sprintf("%d. Between %s and %s", i+1, humanize.IBytes(uint64(r.Min)), humanize.IBytes(uint64(r.Max)))
		heads = append(heads, h)

		count := stat[r].TotalFilesCount
		sz := stat[r].TotalFilesSize

		outputTopStatLine(tw, count, total, sz, h)
	}
	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d file extensions by size:\n\n", Top)
	_, _ = fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

	outputTopTenExtensions(tw, extBySize, total, func(data containers, item *container) (int64, uint64) {
		count := byExt[item.name].Count
		sz := uint64(item.size)
		return count, sz
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d file extensions by count:\n\n", Top)
	_, _ = fmt.Fprintf(tw, format, "Extension", "Count", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "---------", "-----", "------", "----", "------")

	outputTopTenExtensions(tw, extByCount, total, func(data containers, item *container) (int64, uint64) {
		count := item.size
		sz := byExt[item.name].Size
		return count, sz
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d files by size:\n\n", Top)
	_, _ = fmt.Fprintf(tw, "%v\t%v\n", "File", "Size")
	_, _ = fmt.Fprintf(tw, "%v\t%v\n", "------", "----")

	i := 1

	topFiles.Descend(func(c *rbtree.Comparable) bool {
		file := (*c).(*container)
		h := fmt.Sprintf("%d. %s", i, file.name)

		i++

		sz := uint64(file.size)

		_, _ = fmt.Fprintf(tw, "%v\t%v\n", h, humanize.IBytes(sz))

		return true
	})

	_ = tw.Flush()

	_, _ = fmt.Fprintf(w, "\nTOP %d folders by size:\n\n", Top)
	_, _ = fmt.Fprintf(tw, format, "Folder", "Files", "%", "Size", "%")
	_, _ = fmt.Fprintf(tw, format, "------", "-----", "------", "----", "------")

	i = 1

	topFolders.Descend(func(c *rbtree.Comparable) bool {

		folder := (*c).(*container)
		h := fmt.Sprintf("%d. %s", i, folder.name)

		i++

		count := folder.count
		sz := uint64(folder.size)

		outputTopStatLine(tw, count, total, sz, h)

		return true
	})

	_ = tw.Flush()

	if opt.Verbosity && len(opt.Range) > 0 {
		_, _ = fmt.Fprintf(w, "\nDetailed files stat:\n")
		for i, r := range fileSizeRanges {
			if len(filesByRange[r]) == 0 {
				continue
			}

			_, _ = fmt.Fprintf(w, "%s\n", heads[i])

			items := filesByRange[r]
			sort.Sort(sort.Reverse(items))

			for _, item := range items {
				size := humanize.IBytes(uint64(item.size))
				_, _ = fmt.Fprintf(w, "   %s - %s\n", item.name, size)
			}
		}
	}

	printTotals(total, w)
}

func outputTopTenExtensions(tw *tabwriter.Writer, data containers, total totalInfo, selector func(data containers, item *container) (int64, uint64)) {
	for i := 0; i < Top && i < len(data); i++ {
		h := data[i].name

		count, sz := selector(data, data[i])

		outputTopStatLine(tw, count, total, sz, h)
	}
}

func outputTopStatLine(tw *tabwriter.Writer, count int64, total totalInfo, sz uint64, title string) {
	percentOfCount := countPercent(count, total)
	percentOfSize := sizePercent(sz, total)

	_, _ = fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
}

func countPercent(count int64, total totalInfo) float64 {
	return (float64(count) / float64(total.FilesTotal.Count)) * 100
}

func sizePercent(size uint64, total totalInfo) float64 {
	return (float64(size) / float64(total.FilesTotal.Size)) * 100
}

func createSliceFromMap(sizeByExt map[string]countSizeAggregate, mapper func(countSizeAggregate) int64) containers {
	var result = make(containers, len(sizeByExt))
	i := 0
	for k, v := range sizeByExt {
		result[i] = &container{size: mapper(v), name: k}
		i++
	}
	return result
}

func walk(opt options, fs afero.Fs) (totalInfo, map[Range]fileStat, map[Range]containers, map[string]countSizeAggregate, *rbtree.RbTree, *rbtree.RbTree) {
	verboseRanges := make(map[int]bool)
	for _, x := range opt.Range {
		verboseRanges[x] = true
	}
	total := totalInfo{}
	stat := make(map[Range]fileStat)
	filesByRange := make(map[Range]containers)

	byExt := make(map[string]countSizeAggregate)

	start := time.Now()

	filesystemCh := make(chan filesystemItem, 1024)
	go func() {
		walkDirBreadthFirst(opt.Path, fs, filesystemCh)
	}()

	var foldersMu sync.RWMutex
	folders := make(map[string]*container)

	filesChan := make(chan *fileEntry, 1024)

	// Reading filesystem events
	go func() {
		defer close(filesChan)
		for item := range filesystemCh {
			if item.event == fsEventDir {
				foldersMu.Lock()
				folders[item.dir] = &container{name: item.dir}
				total.CountFolders++
				foldersMu.Unlock()
			} else {
				// Only files
				entry := item.entry
				filesChan <- &fileEntry{Size: entry.size, Parent: item.dir, Name: entry.name}
			}
		}
	}()

	topFilesTree := rbtree.NewRbTree()

	// Read all files from channel
	for file := range filesChan {
		fullPath := filepath.Join(file.Parent, file.Name)
		fileContainer := container{size: file.Size, name: fullPath, count: 1}
		updateTopTree(topFilesTree, &fileContainer)

		unsignedSize := uint64(file.Size)

		// Calculate files range statistic
		for i, r := range fileSizeRanges {
			if !r.contains(file.Size) {
				continue
			}

			s := stat[r]
			s.TotalFilesCount++
			s.TotalFilesSize += unsignedSize
			stat[r] = s

			// Store each file info within range only i verbose option set
			if !opt.Verbosity || !verboseRanges[i+1] {
				continue
			}

			nodes, ok := filesByRange[r]
			if !ok {
				filesByRange[r] = make(containers, 0)
			}
			filesByRange[r] = append(nodes, &fileContainer)
		}

		foldersMu.RLock()
		currFolder, ok := folders[file.Parent]

		if ok {
			currFolder.size += file.Size
			currFolder.count++
		}
		foldersMu.RUnlock()

		// Accumulate file statistic
		total.FilesTotal.Count++
		total.FilesTotal.Size += unsignedSize

		ext := filepath.Ext(file.Name)
		a := byExt[ext]
		a.Size += unsignedSize
		a.Count++
		byExt[ext] = a
	}

	topFoldersTree := rbtree.NewRbTree()
	for _, cont := range folders {
		updateTopTree(topFoldersTree, cont)
	}

	total.ReadingTime = time.Since(start)
	return total, stat, filesByRange, byExt, topFoldersTree, topFilesTree
}

func updateTopTree(topTree *rbtree.RbTree, cnt *container) {
	min := topTree.Minimum()
	if topTree.Len() < Top || (*min.Key).(*container).size < cnt.size {
		if topTree.Len() == Top {
			topTree.Delete(min)
		}

		var r rbtree.Comparable
		r = cnt
		node := rbtree.NewNode(&r)
		topTree.Insert(node)
	}
}

func printTotals(t totalInfo, w io.Writer) {

	const totalTemplate = `
Total files:            {{.FilesTotal.Count}} ({{.FilesTotal.Size | toBytesString }})
Total folders:          {{.CountFolders}}
Total file extensions:  {{.CountFileExts}}

Read taken:    {{.ReadingTime}}
`

	var report = template.Must(template.New("totalstat").Funcs(template.FuncMap{"toBytesString": humanize.IBytes}).Parse(totalTemplate))
	_ = report.Execute(w, t)
}