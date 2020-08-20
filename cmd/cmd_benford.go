package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

func newBenford(c conf) *cobra.Command {
	opt := options{}

	var cmd = &cobra.Command{
		Use:     "b",
		Aliases: []string{"benford"},
		Short:   "Show the first digit distribution of files size (benford law validation)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := module.NewContext(top, removeRoot, opt.path)
			benford := module.NewBenfordFileModule(ctx)
			extmod := module.NewExtensionModule(ctx, true)

			run(opt.path, c, benford, extmod)

			return nil
		},
	}

	configure(cmd, &opt)

	return cmd
}
