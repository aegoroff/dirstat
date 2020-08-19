package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

func newFile(c conf) *cobra.Command {
	opt := options{}

	showExtStatistic := false
	showBendfordStatistic := false

	var cmd = &cobra.Command{
		Use:     "fi",
		Aliases: []string{"file"},
		Short:   "Show information about files within folder on volume only",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := module.NewContext(top, removeRoot, opt.path)
			totalmod := module.NewTotalModule(ctx)
			detailfilemod := module.NewDetailFileModule(ctx, opt.vrange)
			totalfilemod := module.NewAggregateFileModule(ctx)
			foldersmod := module.NewFoldersModule(ctx, true)
			extmod := module.NewExtensionModule(ctx, !showExtStatistic)
			benford := module.NewBenfordFileModule(ctx, showBendfordStatistic)

			topfilesmod := module.NewTopFilesModule(ctx)

			run(opt.path, c, totalfilemod, benford, extmod, topfilesmod, detailfilemod, foldersmod, totalmod)

			return nil
		},
	}

	configure(cmd, &opt)

	cmd.Flags().BoolVarP(&showExtStatistic, "ext", "e", false, "Show extensions statistic. By default false")
	cmd.Flags().BoolVarP(&showBendfordStatistic, "bendford", "b", false, "Show files size bendford statistic. By default false")

	return cmd
}
