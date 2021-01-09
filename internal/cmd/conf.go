package cmd

import (
	"github.com/aegoroff/dirstat/internal/out"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type conf interface {
	// fs defines app file system abstraction
	fs() afero.Fs

	env() out.PrintEnvironment

	globals() *globals
}

type appConf struct {
	filesystem afero.Fs
	e          out.PrintEnvironment
	g          *globals
}

type globals struct {
	top        *int
	showMemory *bool
	removeRoot *bool
}

func (a *appConf) fs() afero.Fs {
	return a.filesystem
}

func (a *appConf) env() out.PrintEnvironment {
	return a.e
}

func (a *appConf) globals() *globals {
	return a.g
}

func newAppConf(fs afero.Fs, env out.PrintEnvironment, g *globals) conf {
	c := appConf{
		filesystem: fs,
		e:          env,
		g:          g,
	}
	return &c
}

func confRange(cmd *cobra.Command, rn *[]int) {
	cmd.Flags().IntSliceVarP(rn, "range", "r", []int{}, "Output verbose files info for range specified. Range is the number between 1 and 10")
}
