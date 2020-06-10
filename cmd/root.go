package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
	"runtime"
)

type options struct {
	Verbosity bool
	Range     []int
	Path      string
}

const pathParamName = "path"
const verboseParamName = "verbose"
const rangeParamName = "range"

var top int
var showMemory bool

var appFileSystem afero.Fs
var appWriter io.Writer

func newRoot() *cobra.Command {
	return &cobra.Command{
		Use:   "dirstat",
		Short: "Directory statistic tool",
		Long:  ` A small tool that shows selected folder or drive (on Windows) usage statistic`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	appWriter = os.Stdout
	appFileSystem = afero.NewOsFs()

}

// Execute starts package running
func Execute(args ...string) {
	rootCmd := newRoot()

	rootCmd.PersistentFlags().IntVarP(&top, "top", "t", 10, "The number of lines in top statistics.")
	rootCmd.PersistentFlags().BoolVarP(&showMemory, "memory", "m", false, "Show memory statistic after run")

	rootCmd.AddCommand(newAll())
	rootCmd.AddCommand(newFile())
	rootCmd.AddCommand(newFolder())
	rootCmd.AddCommand(newVersion())

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	printMemUsage(appWriter)
}

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage(w io.Writer) {
	if !showMemory {
		return
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	_, _ = fmt.Fprintf(w, "\nAlloc = %s", humanize.IBytes(m.Alloc))
	_, _ = fmt.Fprintf(w, "\tTotalAlloc = %s", humanize.IBytes(m.TotalAlloc))
	_, _ = fmt.Fprintf(w, "\tSys = %s", humanize.IBytes(m.Sys))
	_, _ = fmt.Fprintf(w, "\tNumGC = %v\n", m.NumGC)
}
