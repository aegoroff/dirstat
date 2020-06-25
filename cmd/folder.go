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
			err := beforeRunCmd(&path, c.fs(), c.w())
			if err != nil {
				return err
			}

			ctx := module.NewContext(top)
			foldersmod := module.NewFoldersModule(ctx, false)
			totalmod := module.NewTotalModule(ctx)
			extmod := module.NewExtensionModule(ctx, true)

			withtiming := module.NewTimeMeasureCommand(module.Execute)
			withtiming(path, c.fs(), c.w(), extmod, foldersmod, totalmod)

			return nil
		},
	}

	configurePath(cmd, &path)

	return cmd
}
