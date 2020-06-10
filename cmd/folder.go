package cmd

import (
	"dirstat/module"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func newFolder() *cobra.Command {
	opt := options{}

	var cmd = &cobra.Command{
		Use:     "fo",
		Aliases: []string{"folder"},
		Short:   "Show information about folders within folder on volume only",
		RunE: func(cmd *cobra.Command, args []string) error {
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

	cmd.Flags().StringVarP(&opt.Path, pathParamName, "p", "", "REQUIRED. Directory path to show info.")

	return cmd
}
