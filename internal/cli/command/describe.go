package command

// import (
// 	"encoding/json"
// 	"fmt"

// 	"github.com/rocketblend/rocketblend/pkg/driver/reference"
// 	"github.com/spf13/cobra"
// )

// // newDescribeCommand creates a new cobra command that fetches the definition of a package.
// // It retrieves the definition based on the reference provided as an argument and formats the output based on the 'output' flag.
// func newDescribeCommand() *cobra.Command {
// 	var output string

// 	cc := &cobra.Command{
// 		Use:   "describe [reference]",
// 		Short: "Fetches a package definition",
// 		Long:  `Fetches the definition of a package by its reference. The output can be formatted by specifying the 'output' flag.`,
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			packages, err := c.factory.GetRocketPackService()
// 			if err != nil {
// 				return fmt.Errorf("failed to get package service: %w", err)
// 			}

// 			ref, err := reference.Parse(args[0])
// 			if err != nil {
// 				return err
// 			}

// 			pkg, err := packages.Get(cmd.Context(), false, ref)
// 			if err != nil {
// 				return err
// 			}

// 			display, err := json.Marshal(pkg)
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
