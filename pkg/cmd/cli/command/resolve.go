package command

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func (srv *Service) newResolveCommand() *cobra.Command {
	var output string

	c := &cobra.Command{
		Use:   "resolve",
		Short: "Output resolved information",
		Long:  `Output the resolved dependencies and paths for the project on the local machine.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			blend, err := srv.findBlendFile()
			if err != nil {
				cmd.Println(err)
				return
			}

			switch output {
			case "json":
				json, err := json.Marshal(blend)
				if err != nil {
					cmd.PrintErrln("failed to resolve config:", err)
					return
				}
				fmt.Println(string(json))
			case "pretty":
				cmd.Println("not implemented")
			default:
				cmd.PrintErrln("invalid output format:", output)
				return
			}
		},
	}

	c.Flags().StringVarP(&output, "output", "o", "pretty", "output format (pretty, json)")

	return c
}
