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
		Short:   "Show information only about folders",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := module.NewContext(*c.globals().top, *c.globals().removeRoot, path)
			foldersmod := module.NewFoldersModule(ctx)
			totalmod := module.NewTotalModule(ctx)

			run(path, c, foldersmod, totalmod)

			return nil
		},
	}

	configurePath(cmd, &path)

	return cmd
}
