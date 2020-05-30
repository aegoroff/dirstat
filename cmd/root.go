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
const top = 10

var appFileSystem = afero.NewOsFs()
var appWriter io.Writer

// rootCmd represents the root command
var rootCmd = &cobra.Command{
	Use:   "dirstat",
	Short: "Directory statistic tool",
	Long:  ` A small tool that shows selected folder or drive (on Windows) usage statistic`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	cobra.MousetrapHelpText = ""
	appWriter = os.Stdout
}

// Execute starts package running
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
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
