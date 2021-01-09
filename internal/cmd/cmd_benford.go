package cmd

import (
	"github.com/aegoroff/dirstat/internal/module"
	"github.com/spf13/cobra"
)

type benfordCmd struct {
	*baseCommand
}

func (b *benfordCmd) execute() error {
	ctx := b.newContext()

	b.run(module.NewBenfordFileModule(ctx, 0), module.NewTotalModule(ctx, 1))

	return nil
}

func newBenford(c conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func(path string) command {
			return &benfordCmd{
				baseCommand: newBaseCmd(c, path),
			}
		},
	}

	cmd := cc.newCobraCommand("b", "benford", "Show the first digit distribution of files size (benford law validation)")

	return cmd
}
