package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
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

	// replaceRoot defines whether to replace root part of full path by <root> macro
	replaceRoot() bool
}

type appConf struct {
	filesystem afero.Fs
	writer     io.Writer
	rr         bool
}

func (a *appConf) replaceRoot() bool {
	return a.rr
}

func (a *appConf) fs() afero.Fs {
	return a.filesystem
}

func (a *appConf) w() io.Writer {
	return a.writer
}

func newAppConf(replaceRoot bool) conf {
	c := appConf{
		filesystem: afero.NewOsFs(),
		writer:     os.Stdout,
		rr:         replaceRoot,
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
	cmd.Flags().StringVarP(path, "path", "p", "", "REQUIRED. Directory path to show info.")
}
