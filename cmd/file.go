package cmd

import (
	"dirstat/module"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const extParamName = "ext"

func newFile() *cobra.Command {
	opt := options{}

	var cmd = &cobra.Command{
		Use:     "fi",
		Aliases: []string{"file"},
		Short:   "Show information about files within folder on volume only",
		RunE: func(cmd *cobra.Command, args []string) error {
			showExtStatistic, err := cmd.Flags().GetBool(extParamName)

			if err != nil {
				return err
			}

			if _, err := appFileSystem.Stat(opt.Path); os.IsNotExist(err) {
				return err
			}

			if opt.Path[len(opt.Path)-1] == ':' {
				opt.Path = filepath.Join(opt.Path, "\\")
			}

			_, _ = fmt.Fprintf(appWriter, "Root: %s\n\n", opt.Path)

			ctx := module.NewContext(top)
			totalmod := module.NewTotalModule(ctx)
			detailfilemod := module.NewDetailFileModule(opt.Verbosity, opt.Range)
			totalfilemod := module.NewTotalFileModule(ctx)
			foldersmod := module.NewFoldersModule(ctx, true)

			var extmod module.Module

			if showExtStatistic {
				extmod = module.NewExtensionModule(ctx, false)
			} else {
				extmod = module.NewExtensionModule(ctx, true)
			}

			topfilesmod := module.NewTopFilesModule(ctx)

			module.Execute(opt.Path, appFileSystem, appWriter, totalfilemod, extmod, topfilesmod, detailfilemod, foldersmod, totalmod)

			return nil
		},
	}

	cmd.Flags().StringVarP(&opt.Path, pathParamName, "p", "", "REQUIRED. Directory path to show info.")
	cmd.Flags().IntSliceVarP(&opt.Range, rangeParamName, "r", []int{}, "Output verbose files info for range specified. Range is the number between 1 and 10")
	cmd.Flags().BoolVarP(&opt.Verbosity, verboseParamName, "v", false, "Be verbose")
	cmd.Flags().BoolP(extParamName, "e", false, "Show extensions statistic. By default false")

	return cmd
}
