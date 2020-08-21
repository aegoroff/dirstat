package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

func newFile(c conf) *cobra.Command {
	opt := options{}

	var cmd = &cobra.Command{
		Use:     "fi",
		Aliases: []string{"file"},
		Short:   "Show information only about files",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := module.NewContext(top, removeRoot, opt.path)
			totalmod := module.NewTotalModule(ctx)
			detailfilemod := module.NewDetailFileModule(ctx, opt.vrange)
			totalfilemod := module.NewAggregateFileModule(ctx)

			topfilesmod := module.NewTopFilesModule(ctx)

			run(opt.path, c, totalfilemod, topfilesmod, detailfilemod, totalmod)

			return nil
		},
	}

	configure(cmd, &opt)

	return cmd
}
