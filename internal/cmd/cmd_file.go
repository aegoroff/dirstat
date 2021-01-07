package cmd

import (
	"github.com/aegoroff/dirstat/internal/module"
	"github.com/spf13/cobra"
)

type fileCmd struct {
	baseCommand
	vrange []int
}

func (f *fileCmd) execute() error {
	ctx := module.NewContext(f.top, f.removeRoot, f.path)

	modules := []module.Module{
		module.NewAggregateFileModule(ctx, 0),
		module.NewTopFilesModule(ctx, 1),
		module.NewDetailFileModule(ctx, 2, f.vrange),
		module.NewTotalModule(ctx, 3),
	}

	f.run(modules...)

	return nil
}

func newFile(c conf) *cobra.Command {
	opt := options{}

	cc := cobraCreator{
		createCmd: func() command {
			return &fileCmd{
				baseCommand: newBaseCmd(c, opt.path),
				vrange:      opt.vrange,
			}
		},
	}

	cmd := cc.newCobraCommand("fi", "file", "Show information only about files")

	configure(cmd, &opt)

	return cmd
}