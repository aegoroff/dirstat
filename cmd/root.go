package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"io"
	"os"
	"runtime"
)

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
}

var showMemory bool
var top int

// Execute starts package running
func Execute(args ...string) {
	rootCmd := newRoot()

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	rootCmd.PersistentFlags().IntVarP(&top, "top", "t", 10, "The number of lines in top statistics.")
	rootCmd.PersistentFlags().BoolVarP(&showMemory, "memory", "m", false, "Show memory statistic after run")

	conf := newAppConf()

	rootCmd.AddCommand(newAll(conf))
	rootCmd.AddCommand(newFile(conf))
	rootCmd.AddCommand(newFolder(conf))
	rootCmd.AddCommand(newVersion(conf.w()))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	printMemUsage(conf.w())
}

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage(w io.Writer) {
	if !showMemory {
		return
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	red := color.FgRed.Render
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	_, _ = fmt.Fprintf(w, "\nAlloc = %s", red(humanize.IBytes(m.Alloc)))
	_, _ = fmt.Fprintf(w, "\tTotalAlloc = %s", red(humanize.IBytes(m.TotalAlloc)))
	_, _ = fmt.Fprintf(w, "\tSys = %s", red(humanize.IBytes(m.Sys)))
	_, _ = fmt.Fprintf(w, "\tNumGC = %v\n", red(m.NumGC))
}
