package command

// import (
// 	"context"

// 	"github.com/rocketblend/rocketblend/pkg/driver/reference"
// 	"github.com/spf13/cobra"
// )

// // newUninstallCommand creates a new cobra.Command that uninstalls dependencies from the current project.
// func newUninstallCommand() *cobra.Command {
// 	cc := &cobra.Command{
// 		Use:   "uninstall [reference]",
// 		Short: "Remove project dependencies",
// 		Long:  "Removes dependencies from the current project.",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			rocketblend, err := c.getDriver()
// 			if err != nil {
// 				return err
// 			}

// 			ref, err := reference.Parse(args[0])
// 			if err != nil {
// 				return err
// 			}

// 			return c.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
// 				return rocketblend.RemoveDependencies(ctx, ref)
// 			}, &spinnerOptions{Suffix: "Removing package..."})
// 		},
// 	}

// 	return cc
// }
