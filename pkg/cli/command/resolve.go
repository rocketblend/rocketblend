package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newResolveCommand creates a new cobra.Command that outputs resolved information about the project.
// This information includes dependencies and paths for the project on the local machine.
func (srv *Service) newResolveCommand() *cobra.Command {
	var output string

	c := &cobra.Command{
		Use:   "resolve",
		Short: "Resolves and outputs project details",
		Long:  `Fetches and prints the resolved dependencies and paths for the project in the local machine in the specified output format (JSON or table)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented")
		},
	}

	c.Flags().StringVarP(&output, "output", "o", "table", "output format (table, json)")

	return c
}
