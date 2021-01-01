package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
)

type options struct {
	vrange []int
	path   string
}

type conf interface {
	// fs defines app file system abstraction
	fs() afero.Fs

	// w defines app output
	w() io.Writer

	globals() *globals
}

type appConf struct {
	filesystem afero.Fs
	writer     io.Writer
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

func (a *appConf) w() io.Writer {
	return a.writer
}

func (a *appConf) globals() *globals {
	return a.g
}

func newAppConf(fs afero.Fs, w io.Writer, g *globals) conf {
	c := appConf{
		filesystem: fs,
		writer:     w,
		g:          g,
	}
	return &c
}

func configure(cmd *cobra.Command, opt *options) {
	configurePath(cmd, &opt.path)
	confRange(cmd, &opt.vrange)
}

func confRange(cmd *cobra.Command, rn *[]int) {
	cmd.Flags().IntSliceVarP(rn, "range", "r", []int{}, "Output verbose files info for range specified. Range is the number between 1 and 10")
}

func configurePath(cmd *cobra.Command, path *string) {
	const param = "path"
	cmd.Flags().StringVarP(path, param, "p", "", "REQUIRED. Directory path to show info.")
	_ = cmd.MarkFlagRequired(param)
}
