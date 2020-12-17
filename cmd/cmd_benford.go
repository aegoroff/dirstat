package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

type benfordCmd struct {
	baseCommand
}

func (b *benfordCmd) execute() error {
	ctx := module.NewContext(b.top, b.removeRoot, b.path)

	var modules []module.Module
	modules = append(modules, module.NewBenfordFileModule(ctx))
	modules = append(modules, module.NewTotalModule(ctx))

	run(b.path, b.c, modules...)

	return nil
}

func newBenford(c conf) *cobra.Command {
	var path string

	cc := cobraCreator{
		createCmd: func() command {
			cmd := benfordCmd{
				baseCommand: newBaseCmd(c, path),
			}
			return &cmd
		},
	}

	cmd := cc.newCobraCommand("b", "benford", "Show the first digit distribution of files size (benford law validation)")

	configurePath(cmd, &path)

	return cmd
}
