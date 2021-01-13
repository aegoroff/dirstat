package cmd

import (
	"github.com/aegoroff/dirstat/internal/module"
	"github.com/spf13/cobra"
)

type allCmd struct {
	*baseCommand
	vrange []int
}

func (a *allCmd) execute() error {
	ctx := a.newContext()

	return a.run(
		module.NewAggregateFileModule(ctx, 0),
		module.NewExtensionModule(ctx, 1),
		module.NewTopFilesModule(ctx, 2),
		module.NewFoldersModule(ctx, 3),
		module.NewDetailFileModule(ctx, 4, a.vrange),
		module.NewTotalModule(ctx, 5),
	)
}

func newAll(c conf) *cobra.Command {
	var vrange []int

	cc := cobraCreator{
		createCmd: func(path string) command {
			return &allCmd{
				baseCommand: newBaseCmd(c, path),
				vrange:      vrange,
			}
		},
	}

	cmd := cc.newCobraCommand("a", "all", "Show all information about folder/volume")

	confRange(cmd, &vrange)

	return cmd
}
