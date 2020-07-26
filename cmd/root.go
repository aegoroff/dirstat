package cmd

import (
	"github.com/spf13/cobra"
	"os"
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
var removeRoot bool

// Execute starts package running
func Execute(args ...string) {
	rootCmd := newRoot()

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	rootCmd.PersistentFlags().IntVarP(&top, "top", "t", 10, "The number of lines in top statistics.")
	rootCmd.PersistentFlags().BoolVarP(&showMemory, "memory", "m", false, "Show memory statistic after run")
	rootCmd.PersistentFlags().BoolVarP(&removeRoot, "removeroot", "o", false, "Remove root part from full path i.e. output relative paths")

	conf := newAppConf()

	rootCmd.AddCommand(newAll(conf))
	rootCmd.AddCommand(newFile(conf))
	rootCmd.AddCommand(newFolder(conf))
	rootCmd.AddCommand(newVersion(conf.w()))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
