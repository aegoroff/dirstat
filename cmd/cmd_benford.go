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

	b.run(module.NewBenfordFileModule(ctx, 0), module.NewTotalModule(ctx, 1))

	return nil
}

func newBenford(c conf) *cobra.Command {
	var path string

	cc := cobraCreator{
		createCmd: func() command {
			return &benfordCmd{
				baseCommand: newBaseCmd(c, path),
			}
		},
	}

	cmd := cc.newCobraCommand("b", "benford", "Show the first digit distribution of files size (benford law validation)")

	configurePath(cmd, &path)

	return cmd
}
