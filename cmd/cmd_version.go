package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.6.0"

func newVersion(w io.Writer) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "Print the version number of dirstat",
		Long:    `All software has versions. This is dirstat's`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(w, "dirstat v%s\n", Version)
			return err
		},
	}
	return cmd
}
