package cmd

import (
	"github.com/aegoroff/dirstat/internal/module"
	"github.com/spf13/cobra"
)

type extCmd struct {
	*baseCommand
}

func (e *extCmd) execute() error {
	ctx := e.newContext()

	e.run(module.NewExtensionModule(ctx, 0), module.NewTotalModule(ctx, 1))

	return nil
}

func newExt(c conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func(path string) command {
			return &extCmd{
				baseCommand: newBaseCmd(c, path),
			}
		},
	}

	cmd := cc.newCobraCommand("e", "ext", "Show file extensions statistic")

	return cmd
}
