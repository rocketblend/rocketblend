package command

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// newRunCommand creates a new cobra command for running the project.
// It sets up all necessary flags and executes the project through the driver.
func (srv *Service) newRunCommand() *cobra.Command {
	var background bool

	c := &cobra.Command{
		Use:   "run",
		Short: "Runs the project",
		Long:  `Launches the project in the current working directory. Can optionally run in the background.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			blend, err := srv.findBlendFilePath(srv.flags.workingDirectory)
			if err != nil {
				return fmt.Errorf("unable to locate project file: %w", err)
			}

			err = srv.run(cmd.Context(), blend, background)
			if err != nil {
				return fmt.Errorf("failed to run project: %w", err)
			}

			return nil
		},
	}

	// background flag allows the project to be run in the background.
	c.Flags().BoolVarP(&background, "background", "b", false, "run the project in the background")

	return c
}

func (srv *Service) run(ctx context.Context, blendPath string, background bool) error {
	rocketblend, err := srv.factory.CreateRocketBlendService()
	if err != nil {
		return fmt.Errorf("failed to create rocketblend: %w", err)
	}

	blendFile, err := rocketblend.Load(srv.flags.workingDirectory)
	if err != nil {
		return fmt.Errorf("failed to load blend file: %w", err)
	}

	cmd, err := rocketblend.GetCMD(ctx, blendFile, background, []string{})
	if err != nil {
		return err
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open: %s", err)
	}

	return nil
}
