package cmd

import (
	"github.com/aegoroff/dirstat/internal/module"
	"github.com/spf13/cobra"
)

type fileCmd struct {
	*baseCommand
	vrange []int
}

func (f *fileCmd) execute() error {
	ctx := f.newContext()

	return f.run(
		module.NewAggregateFileModule(ctx, 0),
		module.NewTopFilesModule(ctx, 1),
		module.NewDetailFileModule(ctx, 2, f.vrange),
		module.NewTotalModule(ctx, 3),
	)
}

func newFile(c conf) *cobra.Command {
	var vrange []int

	cc := cobraCreator{
		createCmd: func(path string) command {
			return &fileCmd{
				baseCommand: newBaseCmd(c, path),
				vrange:      vrange,
			}
		},
	}

	cmd := cc.newCobraCommand("fi", "file", "Show information only about files")

	confRange(cmd, &vrange)

	return cmd
}
