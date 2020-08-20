package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

func newFile(c conf) *cobra.Command {
	opt := options{}

	showExtStatistic := false

	var cmd = &cobra.Command{
		Use:     "fi",
		Aliases: []string{"file"},
		Short:   "Show information about files within folder on volume only",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := module.NewContext(top, removeRoot, opt.path)
			totalmod := module.NewTotalModule(ctx)
			detailfilemod := module.NewDetailFileModule(ctx, opt.vrange)
			totalfilemod := module.NewAggregateFileModule(ctx)
			var extmod module.Module

			if showExtStatistic {
				extmod = module.NewExtensionModule(ctx)
			} else {
				extmod = module.NewVoidModule()
			}

			topfilesmod := module.NewTopFilesModule(ctx)

			run(opt.path, c, totalfilemod, extmod, topfilesmod, detailfilemod, totalmod)

			return nil
		},
	}

	configure(cmd, &opt)

	cmd.Flags().BoolVarP(&showExtStatistic, "ext", "e", false, "Show extensions statistic. By default false")

	return cmd
}
