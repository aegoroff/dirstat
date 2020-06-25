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
			err := beforeRunCmd(&opt.path, c.fs(), c.w())
			if err != nil {
				return err
			}

			ctx := module.NewContext(top)
			totalmod := module.NewTotalModule(ctx)
			detailfilemod := module.NewDetailFileModule(opt.verbosity, opt.vrange)
			totalfilemod := module.NewAggregateFileModule(ctx)
			foldersmod := module.NewFoldersModule(ctx, true)
			extmod := module.NewExtensionModule(ctx, !showExtStatistic)

			topfilesmod := module.NewTopFilesModule(ctx)

			withtiming := module.NewTimeMeasureCommand(module.Execute)
			withtiming(opt.path, c.fs(), c.w(), totalfilemod, extmod, topfilesmod, detailfilemod, foldersmod, totalmod)

			return nil
		},
	}

	configure(cmd, &opt)

	cmd.Flags().BoolVarP(&showExtStatistic, "ext", "e", false, "Show extensions statistic. By default false")

	return cmd
}
