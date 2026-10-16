package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aegoroff/dirstat/internal/module"
	"github.com/aegoroff/dirstat/scan"
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

func newBaseCmd(c conf, path string) *baseCommand {
	return &baseCommand{
		c:          c,
		top:        *c.globals().top,
		removeRoot: *c.globals().removeRoot,
		path:       correctPath(path),
	}
}

// correctPath turns Windows drive designator like "C:" into drive root "C:\"
func correctPath(path string) string {
	if len(path) > 0 && path[len(path)-1] == ':' {
		return filepath.Join(path, string(os.PathSeparator))
	}
	return path
}

func (b *baseCommand) run(modules ...module.Module) error {
	if b.path == "" {
		return fmt.Errorf("path must not be empty")
	}

	c := b.c
	if _, err := c.fs().Stat(b.path); err != nil {
		return err
	}

	var r runner
	{
		r = module.Execute
		r = newTimeMeasureR(r)
		r = newPrintRootR(r)
	}

	p, err := c.env().NewPrinter()
	if err != nil {
		return err
	}

	defer scan.Close(c.env().Writer())

	if *c.globals().showMemory {
		r = newPrintMemoryR(r)
	}

	r(b.path, c.fs(), p, modules...)
	return nil
}

func (b *baseCommand) newContext() *module.Context {
	return module.NewContext(b.top, b.removeRoot, b.path)
}

type cobraCreator struct {
	createCmd func(path string) command
}

func (c *cobraCreator) runE() cobraRunSignature {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			// Errors from here on are runtime ones, not usage mistakes
			cmd.SilenceUsage = true
			return c.createCmd(args[0]).execute()
		}
		return cmd.Help()
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
