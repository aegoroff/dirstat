package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

func newFolder(c conf) *cobra.Command {
	var path string

	var cmd = &cobra.Command{
		Use:     "fo",
		Aliases: []string{"folder"},
		Short:   "Show information about folders within folder on volume only",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := module.NewContext(top)
			foldersmod := module.NewFoldersModule(ctx, false)
			totalmod := module.NewTotalModule(ctx)
			extmod := module.NewExtensionModule(ctx, true)

			run(path, c, extmod, foldersmod, totalmod)

			return nil
		},
	}

	configurePath(cmd, &path)

	return cmd
}
