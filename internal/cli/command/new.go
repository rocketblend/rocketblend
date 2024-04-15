package command

// import (
// 	"context"
// 	"fmt"
// 	"path/filepath"
// 	"strings"

// 	"github.com/rocketblend/rocketblend/pkg/driver/blendconfig"
// 	"github.com/rocketblend/rocketblend/pkg/driver/rocketfile"
// 	"github.com/spf13/cobra"
// )

// // newNewCommand creates a new cobra.Command object initialized for creating a new project.
// // It expects a single argument which is the name of the project.
// // It uses the 'skip-install' flag to decide whether or not to install dependencies.
// func (c *cli) newNewCommand() *cobra.Command {
// 	var skipInstall bool

// 	cc := &cobra.Command{
// 		Use:   "new [name]",
// 		Short: "Create a new project",
// 		Long:  `Creates a new project with a specified name. Use the 'skip-install' flag to skip installing dependencies.`,
// 		Args:  cobra.MinimumNArgs(1),
// 		PreRunE: func(cmd *cobra.Command, args []string) error {
// 			return c.validateProjectName(args[0])
// 		},
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			config, err := c.getConfig()
// 			if err != nil {
// 				return err
// 			}

// 			blendConfig, err := blendconfig.New(
// 				c.flags.workingDirectory,
// 				c.ensureBlendExtension(args[0]),
// 				rocketfile.New(config.DefaultBuild),
// 			)
// 			if err != nil {
// 				return err
// 			}

// 			driver, err := c.createDriver(blendConfig)
// 			if err != nil {
// 				return err
// 			}

// 			return c.runWithSpinner(cmd.Context(), func(ctx context.Context) error {
// 				return driver.Create(ctx)
// 			}, &spinnerOptions{Suffix: "Creating project..."})
// 		},
// 	}

// 	cc.Flags().BoolVarP(&skipInstall, "skip-install", "s", false, "skip installing dependencies")

// 	return cc
// }

// // validateProjectName checks if the project name is valid.
// func (c *cli) validateProjectName(projectName string) error {
// 	if filepath.IsAbs(projectName) || strings.Contains(projectName, string(filepath.Separator)) {
// 		return fmt.Errorf("%q is not a valid project name, it should not contain any path separators", projectName)
// 	}

// 	if ext := filepath.Ext(projectName); ext != "" {
// 		return fmt.Errorf("%q is not a valid project name, it should not contain any file extension", projectName)
// 	}

// 	return nil
// }

// // ensureBlendExtension adds ".blend" extension to filename if it does not already have it.
// func (c *cli) ensureBlendExtension(filename string) string {
// 	if !strings.HasSuffix(filename, ".blend") {
// 		filename += ".blend"
// 	}

// 	return filename
// }
