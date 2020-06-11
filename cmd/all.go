package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

func newAll(c conf) *cobra.Command {
	opt := options{}

	var cmd = &cobra.Command{
		Use:     "a",
		Aliases: []string{"all"},
		Short:   "Show all information about folder/volume",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := beforeRunCmd(&opt.path, c.fs(), c.w())
			if err != nil {
				return err
			}

			ctx := module.NewContext(top)
			foldersmod := module.NewFoldersModule(ctx, false)
			totalmod := module.NewTotalModule(ctx)
			detailfilemod := module.NewDetailFileModule(opt.verbosity, opt.vrange)
			totalfilemod := module.NewTotalFileModule(ctx)
			extmod := module.NewExtensionModule(ctx, false)
			topfilesmod := module.NewTopFilesModule(ctx)

			module.Execute(opt.path, c.fs(), c.w(), totalfilemod, extmod, topfilesmod, foldersmod, detailfilemod, totalmod)

			return nil
		},
	}

	configure(cmd, &opt)

	return cmd
}
