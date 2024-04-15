package command

// import (
// 	"encoding/json"

// 	"github.com/spf13/cobra"
// )

// // newResolveCommand creates a new cobra.Command that outputs resolved information about the project.
// // This information includes dependencies and paths for the project on the local machine.
// func newResolveCommand() *cobra.Command {
// 	var output string

// 	cc := &cobra.Command{
// 		Use:   "resolve",
// 		Short: "Resolves and outputs project details",
// 		Long:  `Fetches and prints the resolved dependencies and paths for the project in the local machine in the specified output format (JSON or table)`,
// 		Args:  cobra.NoArgs,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			rocketblend, err := c.getDriver()
// 			if err != nil {
// 				return err
// 			}

// 			blendFile, err := rocketblend.ResolveBlendFile(cmd.Context())
// 			if err != nil {
// 				return err
// 			}

// 			display, err := json.Marshal(blendFile)
// 			if err != nil {
// 				return err
// 			}

// 			cmd.Println(string(display))

// 			return nil
// 		},
// 	}

// 	cc.Flags().StringVarP(&output, "output", "o", "table", "output format (table, json)")

// 	return cc
// }
