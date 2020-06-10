package cmd

import (
	"dirstat/module"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func newAll() *cobra.Command {
	opt := options{}

	var cmd = &cobra.Command{
		Use:     "a",
		Aliases: []string{"all"},
		Short:   "Show all information about folder/volume",
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
			detailfilemod := module.NewDetailFileModule(opt.Verbosity, opt.Range)
			totalfilemod := module.NewTotalFileModule(ctx)
			extmod := module.NewExtensionModule(ctx, false)
			topfilesmod := module.NewTopFilesModule(ctx)

			module.Execute(opt.Path, appFileSystem, appWriter, totalfilemod, extmod, topfilesmod, foldersmod, detailfilemod, totalmod)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opt.Path, pathParamName, "p", "", "REQUIRED. Directory path to show info.")
	cmd.Flags().IntSliceVarP(&opt.Range, rangeParamName, "r", []int{}, "Output verbose files info for range specified. Range is the number between 1 and 10")
	cmd.Flags().BoolVarP(&opt.Verbosity, verboseParamName, "v", false, "Be verbose")

	return cmd
}
