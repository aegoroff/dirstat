package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
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

// Execute starts package running
func Execute(fs afero.Fs, w io.Writer, args ...string) error {
	rootCmd := newRoot()

	if args != nil && len(args) > 0 {
		rootCmd.SetArgs(args)
	}

	var top int
	var showMemory bool
	var removeRoot bool

	rootCmd.PersistentFlags().IntVarP(&top, "top", "t", 10, "The number of lines in top statistics.")
	rootCmd.PersistentFlags().BoolVarP(&showMemory, "memory", "m", false, "Show memory statistic after run")
	rootCmd.PersistentFlags().BoolVarP(&removeRoot, "removeroot", "o", false, "Remove root part from full path i.e. output relative paths")

	g := globals{
		top:        &top,
		showMemory: &showMemory,
		removeRoot: &removeRoot,
	}
	conf := newAppConf(fs, w, &g)

	rootCmd.AddCommand(newAll(conf))
	rootCmd.AddCommand(newFile(conf))
	rootCmd.AddCommand(newFolder(conf))
	rootCmd.AddCommand(newBenford(conf))
	rootCmd.AddCommand(newExt(conf))
	rootCmd.AddCommand(newVersion(conf.w()))

	return rootCmd.Execute()
}