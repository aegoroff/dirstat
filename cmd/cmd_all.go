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

	var modules []module.Module
	modules = append(modules, module.NewAggregateFileModule(ctx))
	modules = append(modules, module.NewBenfordFileModule(ctx))
	modules = append(modules, module.NewExtensionModule(ctx))
	modules = append(modules, module.NewTopFilesModule(ctx))
	modules = append(modules, module.NewFoldersModule(ctx))
	modules = append(modules, module.NewDetailFileModule(ctx, a.vrange))
	modules = append(modules, module.NewTotalModule(ctx))

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
