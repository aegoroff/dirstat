package cmd

import (
	"github.com/aegoroff/dirstat/internal/module"
	"github.com/spf13/cobra"
)

type folderCmd struct {
	*baseCommand
}

func (f *folderCmd) execute() error {
	ctx := f.newContext()

	return f.run(module.NewFoldersModule(ctx, 0), module.NewTotalModule(ctx, 1))
}

func newFolder(c conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func(path string) command {
			return &folderCmd{
				baseCommand: newBaseCmd(c, path),
			}
		},
	}

	cmd := cc.newCobraCommand("fo", "folder", "Show information only about folders")

	return cmd
}
