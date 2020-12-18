package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

type extCmd struct {
	baseCommand
}

func (e *extCmd) execute() error {
	ctx := module.NewContext(e.top, e.removeRoot, e.path)

	run(e.path, e.c, module.NewExtensionModule(ctx), module.NewTotalModule(ctx))

	return nil
}

func newExt(c conf) *cobra.Command {
	var path string

	cc := cobraCreator{
		createCmd: func() command {
			cmd := extCmd{
				baseCommand: newBaseCmd(c, path),
			}
			return &cmd
		},
	}

	cmd := cc.newCobraCommand("e", "ext", "Show file extensions statistic")

	configurePath(cmd, &path)

	return cmd
}
