package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"text/tabwriter"
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

		execute(opt.Path, appFileSystem, appWriter, opt.Verbosity, opt.Range)

		printMemUsage(appWriter)
		return nil
	},
}

func init() {
	cobra.MousetrapHelpText = ""
	appWriter = os.Stdout
	rootCmd.Flags().StringP(pathParamName, "p", "", "REQUIRED. Directory path to show info.")
	rootCmd.Flags().IntSliceP(rangeParamName, "r", []int{}, "Output verbose files info for range specified. Range is the number between 1 and 10")
	rootCmd.Flags().BoolP(verboseParamName, "v", false, "Be verbose")
}

// Execute starts package running
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func outputTopStatLine(tw *tabwriter.Writer, count int64, total *totalInfo, sz uint64, title string) {
	percentOfCount := countPercent(count, total)
	percentOfSize := sizePercent(sz, total)

	_, _ = fmt.Fprintf(tw, "%v\t%v\t%.2f%%\t%v\t%.2f%%\n", title, count, percentOfCount, humanize.IBytes(sz), percentOfSize)
}

func countPercent(count int64, total *totalInfo) float64 {
	return (float64(count) / float64(total.FilesTotal.Count)) * 100
}

func sizePercent(size uint64, total *totalInfo) float64 {
	return (float64(size) / float64(total.FilesTotal.Size)) * 100
}

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage(w io.Writer) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	_, _ = fmt.Fprintf(w, "\nAlloc = %s", humanize.IBytes(m.Alloc))
	_, _ = fmt.Fprintf(w, "\tTotalAlloc = %s", humanize.IBytes(m.TotalAlloc))
	_, _ = fmt.Fprintf(w, "\tSys = %s", humanize.IBytes(m.Sys))
	_, _ = fmt.Fprintf(w, "\tNumGC = %v\n", m.NumGC)
}
