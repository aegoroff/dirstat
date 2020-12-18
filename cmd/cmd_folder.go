package cmd

import (
	"dirstat/module"
	"github.com/spf13/cobra"
)

type folderCmd struct {
	baseCommand
}

func (f *folderCmd) execute() error {
	ctx := module.NewContext(f.top, f.removeRoot, f.path)

	run(f.path, f.c, module.NewFoldersModule(ctx), module.NewTotalModule(ctx))

	return nil
}

func newFolder(c conf) *cobra.Command {
	var path string

	cc := cobraCreator{
		createCmd: func() command {
			return &folderCmd{
				baseCommand: newBaseCmd(c, path),
			}
		},
	}

	cmd := cc.newCobraCommand("fo", "folder", "Show information only about folders")

	configurePath(cmd, &path)

	return cmd
}
