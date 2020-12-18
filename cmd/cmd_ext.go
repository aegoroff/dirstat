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

	e.run(module.NewExtensionModule(ctx), module.NewTotalModule(ctx))

	return nil
}

func newExt(c conf) *cobra.Command {
	var path string

	cc := cobraCreator{
		createCmd: func() command {
			return &extCmd{
				baseCommand: newBaseCmd(c, path),
			}
		},
	}

	cmd := cc.newCobraCommand("e", "ext", "Show file extensions statistic")

	configurePath(cmd, &path)

	return cmd
}
