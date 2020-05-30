package cmd

import (
	"dirstat/module"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:     "a",
	Aliases: []string{"all"},
	Short:   "Show all information about folder/volume",
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

		opt := options{Verbosity: verbose, Path: path, Range: ranges}

		if _, err := appFileSystem.Stat(opt.Path); os.IsNotExist(err) {
			return err
		}

		if opt.Path[len(opt.Path)-1] == ':' {
			opt.Path = filepath.Join(opt.Path, "\\")
		}

		_, _ = fmt.Fprintf(appWriter, "Root: %s\n\n", opt.Path)

		ctx := module.NewContext(top)
		foldersmod := module.NewFoldersModule(ctx)
		totalmod := module.NewTotalModule(ctx)
		rangemod := module.NewRangeModule(ctx, opt.Verbosity, opt.Range)
		totalfilemod := module.NewTotalFileModule(ctx)
		extmod := module.NewExtensionModule(ctx)
		topfilesmod := module.NewTopFilesModule(ctx)

		modules := []module.Module{totalfilemod, extmod, topfilesmod, foldersmod, rangemod, totalmod}

		module.Execute(opt.Path, appFileSystem, appWriter, modules)

		printMemUsage(appWriter)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().StringP(pathParamName, "p", "", "REQUIRED. Directory path to show info.")
	allCmd.Flags().IntSliceP(rangeParamName, "r", []int{}, "Output verbose files info for range specified. Range is the number between 1 and 10")
	allCmd.Flags().BoolP(verboseParamName, "v", false, "Be verbose")
}
