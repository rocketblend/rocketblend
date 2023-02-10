package command

import "github.com/spf13/cobra"

func (srv *Service) newResolveCommand() *cobra.Command {
	var output string

	c := &cobra.Command{
		Use:   "resolve",
		Short: "Output resolved information",
		Long:  `Output the resolved dependencies and paths for the project on the local machine.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("resolve called")
		},
	}

	c.Flags().StringVarP(&output, "output", "o", "pretty", "output format (pretty, json)")

	return c
}
