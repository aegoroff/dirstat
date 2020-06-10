package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.3.0"

func newVersion() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"ver"},
		Short:   "Print the version number of dirstat",
		Long:    `All software has versions. This is dirstat's`,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintf(appWriter, "dirstat v%s\n", Version)
			return err
		},
	}
	return cmd
}
