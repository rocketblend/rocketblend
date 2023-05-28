package command

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// newStartCommand creates a new cobra.Command that starts the project located in the working directory.
// It optionally allows the project to be started in the background.
func (srv *Service) newStartCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "start",
		Short: "Starts the project",
		Long:  `Starts the project located in the current working directory.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			blendPath, err := srv.findBlendFilePath(srv.flags.workingDirectory)
			if err != nil {
				return fmt.Errorf("unable to locate project file: %w", err)
			}

			err = srv.start(cmd.Context(), blendPath, []string{})
			if err != nil {
				return fmt.Errorf("failed to start project: %w", err)
			}

			return nil
		},
	}

	return c
}

func (srv *Service) start(ctx context.Context, blendPath string, args []string) error {
	rocketblend, err := srv.factory.CreateRocketBlendService()
	if err != nil {
		return fmt.Errorf("failed to create rocketblend: %w", err)
	}

	blendFile, err := rocketblend.Load(srv.flags.workingDirectory)
	if err != nil {
		return fmt.Errorf("failed to load blend file: %w", err)
	}

	cmd, err := rocketblend.GetCMD(ctx, blendFile, false, args)
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open: %s", err)
	}

	return nil
}
