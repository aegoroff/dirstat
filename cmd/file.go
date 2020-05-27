package cmd

import (
	"dirstat/module"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const extParamName = "ext"

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:     "fi",
	Aliases: []string{"file"},
	Short:   "Show information about files within folder on volume only",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := cmd.Flags().GetString(pathParamName)

		if err != nil {
			return err
		}

		ranges, err := cmd.Flags().GetIntSlice(rangeParamName)

		if err != nil {
			return err
		}

		verbose, err := cmd.Flags().GetBool(verboseParamName)

		if err != nil {
			return err
		}

		showExtStatistic, err := cmd.Flags().GetBool(extParamName)

		if err != nil {
			return err
		}

		opt := options{Verbosity: verbose, Path: path, Range: ranges}

		if _, err := appFileSystem.Stat(opt.Path); os.IsNotExist(err) {
			return err
		}

		if opt.Path[len(opt.Path)-1] == ':' {
			opt.Path = filepath.Join(opt.Path, "\\")
		}

		_, _ = fmt.Fprintf(appWriter, "Root: %s\n\n", opt.Path)

		ctx := module.NewContext()
		totalmod := module.NewTotalModule(ctx)
		rangemod := module.NewRangeModule(ctx, opt.Verbosity, opt.Range)
		totalfilemod := module.NewTotalFileModule(ctx)
		foldersmod := module.NewFoldersHiddenModule(ctx)

		var extmod module.Module

		if showExtStatistic {
			extmod = module.NewExtensionModule(ctx)
		} else {
			extmod = module.NewExtensionHiddenModule(ctx)
		}

		topfilesmod := module.NewTopFilesModule(ctx)

		modules := []module.Module{totalfilemod, extmod, topfilesmod, rangemod, foldersmod, totalmod}

		module.Execute(opt.Path, appFileSystem, appWriter, modules)

		printMemUsage(appWriter)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
	fileCmd.Flags().StringP(pathParamName, "p", "", "REQUIRED. Directory path to show info.")
	fileCmd.Flags().IntSliceP(rangeParamName, "r", []int{}, "Output verbose files info for range specified. Range is the number between 1 and 10")
	fileCmd.Flags().BoolP(verboseParamName, "v", false, "Be verbose")
	fileCmd.Flags().BoolP(extParamName, "e", false, "Show extensions statistic. By default false")
}
