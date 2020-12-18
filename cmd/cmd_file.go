package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

type fileCmd struct {
	baseCommand
	vrange []int
}

func (f *fileCmd) execute() error {
	ctx := module.NewContext(f.top, f.removeRoot, f.path)

	var modules []module.Module
	modules = append(modules, module.NewAggregateFileModule(ctx))
	modules = append(modules, module.NewTopFilesModule(ctx))
	modules = append(modules, module.NewDetailFileModule(ctx, f.vrange))
	modules = append(modules, module.NewTotalModule(ctx))

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
