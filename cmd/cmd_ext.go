package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

func newExt(c conf) *cobra.Command {
	opt := options{}

	var cmd = &cobra.Command{
		Use:     "e",
		Aliases: []string{"ext"},
		Short:   "Show file extensions statistic",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := module.NewContext(top, removeRoot, opt.path)
			totalmod := module.NewTotalModule(ctx)
			extmod := module.NewExtensionModule(ctx)

			run(opt.path, c, extmod, totalmod)

			return nil
		},
	}

	configure(cmd, &opt)

	return cmd
}
