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
	Short:   "Show information about folders within folder on volume only",
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

		ctx := module.NewContext(top)
		foldersmod := module.NewFoldersModule(ctx, false)
		totalmod := module.NewTotalModule(ctx)
		extmod := module.NewExtensionModule(ctx, true)

		module.Execute(opt.Path, appFileSystem, appWriter, extmod, foldersmod, totalmod)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)
	folderCmd.Flags().StringP(pathParamName, "p", "", "REQUIRED. Directory path to show info.")
}
