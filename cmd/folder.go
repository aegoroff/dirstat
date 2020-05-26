package cmd

import (
	"dirstat/module"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// folderCmd represents the folder command
var folderCmd = &cobra.Command{
	Use:     "fo",
	Aliases: []string{"folder"},
	Short:   "Show information about folders only",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString(pathParamName)

		if err != nil {
			return err
		}

		opt := options{Path: path}

		if _, err := appFileSystem.Stat(opt.Path); os.IsNotExist(err) {
			return err
		}

		if opt.Path[len(opt.Path)-1] == ':' {
			opt.Path = filepath.Join(opt.Path, "\\")
		}

		_, _ = fmt.Fprintf(appWriter, "Root: %s\n\n", opt.Path)

		ctx := module.NewContext()
		foldersmod := module.NewFoldersModule(ctx)
		totalmod := module.NewTotalModule(ctx)
		rangemod := module.NewRangeHiddenModule(ctx)
		totalfilemod := module.NewTotalHiddenFileModule(ctx)
		extmod := module.NewExtensionHiddenModule(ctx)
		topfilesmod := module.NewTopFilesHiddenModule(ctx)

		modules := []module.Module{totalfilemod, extmod, topfilesmod, foldersmod, rangemod, totalmod}

		module.Execute(opt.Path, appFileSystem, appWriter, ctx, modules)

		printMemUsage(appWriter)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)
	folderCmd.Flags().StringP(pathParamName, "p", "", "REQUIRED. Directory path to show info.")
}
