package command

// import (
// 	"context"

// 	"github.com/rocketblend/rocketblend/pkg/driver/reference"
// 	"github.com/spf13/cobra"
// )

// // newInstallCommand creates a new cobra command for installing project dependencies.
// func newInstallCommand() *cobra.Command {
// 	var forceUpdate bool

// 	cc := &cobra.Command{
// 		Use:   "install [reference]",
// 		Short: "Installs project dependencies",
// 		Long:  `Adds the specified dependencies to the current project and installs them. If no reference is provided, all dependencies in the project are installed instead.`,
// 		Args:  cobra.MaximumNArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			rocketblend, err := c.getDriver()
// 			if err != nil {
// 				return err
// 			}

// 			if len(args) > 0 {
// 				ref, err := reference.Parse(args[0])
// 				if err != nil {
// 					return err
// 				}

// 				// Add and installs the dependency to the project.
// 				return c.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
// 					return rocketblend.AddDependencies(ctx, forceUpdate, ref)
// 				}, &spinnerOptions{Suffix: "Installing package..."})
// 			}

// 			// Installs all dependencies in the project.
// 			return c.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
// 				return rocketblend.InstallDependencies(ctx)
// 			}, &spinnerOptions{Suffix: "Installing dependencies..."})
// 		},
// 	}

// 	cc.Flags().BoolVarP(&forceUpdate, "update", "u", false, "refreshes the package definition before installing it")

// 	return cc
// }
