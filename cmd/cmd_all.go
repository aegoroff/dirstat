package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

type allCmd struct {
	baseCommand
	vrange []int
}

func (a *allCmd) execute() error {
	ctx := module.NewContext(a.top, a.removeRoot, a.path)

	modules := []module.Module{
		module.NewAggregateFileModule(ctx, 0),
		module.NewBenfordFileModule(ctx, 1),
		module.NewExtensionModule(ctx, 2),
		module.NewTopFilesModule(ctx, 3),
		module.NewFoldersModule(ctx, 4),
		module.NewDetailFileModule(ctx, 5, a.vrange),
		module.NewTotalModule(ctx, 6),
	}

	a.run(modules...)

	return nil
}

func newAll(c conf) *cobra.Command {
	opt := options{}

	cc := cobraCreator{
		createCmd: func() command {
			return &allCmd{
				baseCommand: newBaseCmd(c, opt.path),
				vrange:      opt.vrange,
			}
		},
	}

	cmd := cc.newCobraCommand("a", "all", "Show all information about folder/volume")

	configure(cmd, &opt)

	return cmd
}
