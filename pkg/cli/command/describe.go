package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newDescribeCommand creates a new cobra command that fetches the definition of a package.
// It retrieves the definition based on the reference provided as an argument and formats the output based on the 'output' flag.
func (srv *Service) newDescribeCommand() *cobra.Command {
	var output string

	c := &cobra.Command{
		Use:   "describe [reference]",
		Short: "Fetches a package definition",
		Long:  `Fetches the definition of a package by its reference. The output can be formatted by specifying the 'output' flag.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented")
		},
	}

	c.Flags().StringVarP(&output, "output", "o", "table", "output format (table, json)")

	return c
}
