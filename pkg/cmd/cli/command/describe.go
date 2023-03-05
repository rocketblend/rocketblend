package command

import (
	"encoding/json"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newDescribeCommand() *cobra.Command {
	var output string

	c := &cobra.Command{
		Use:   "describe [reference]",
		Short: "Provide a description for a rocketpack",
		Long:  `Provide a description for a rocketpack`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			reference, err := reference.Parse(args[0])
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			pack, err := srv.driver.DescribePackByReference(reference)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			json, err := json.Marshal(pack)
			if err != nil {
				cmd.PrintErrln("failed to describe pack:", err)
				return
			}

			cmd.Println(string(json))
		},
	}

	c.Flags().StringVarP(&output, "output", "o", "json", "output format")

	return c
}
