package cmd

import (
	"github.com/aegoroff/dirstat/internal/module"
	"github.com/spf13/cobra"
)

type command interface {
	execute() error
}

type cobraRunSignature func(cmd *cobra.Command, args []string) error

type baseCommand struct {
	c          conf
	path       string
	top        int
	removeRoot bool
}

func newBaseCmd(c conf, path string) baseCommand {
	return baseCommand{
		c:          c,
		top:        *c.globals().top,
		removeRoot: *c.globals().removeRoot,
		path:       path,
	}
}

func (b *baseCommand) run(modules ...module.Module) {
	var r runner
	{
		r = module.Execute
		r = newTimeMeasureR(r)
		r = newPathCorrectionR(r)
	}

	c := b.c

	if *c.globals().showMemory {
		r = newPrintMemoryR(r)
	}

	r(b.path, c.fs(), c.w(), modules...)
}

type cobraCreator struct {
	createCmd func(path string) command
}

func (c *cobraCreator) runE() cobraRunSignature {
	return func(cmd *cobra.Command, args []string) error {
		path := ""
		if len(args) > 0 {
			path = args[0]
		} else {
			return cmd.Help()
		}
		return c.createCmd(path).execute()
	}
}

func (c *cobraCreator) newCobraCommand(use, alias, short string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     use + " [path]",
		Aliases: []string{alias},
		Short:   short,
		Args:    cobra.MaximumNArgs(1),
		RunE:    c.runE(),
	}
	return cmd
}
