package command

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/rocketblend"
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
			blend, err := srv.findBlendFile(srv.flags.workingDirectory)
			if err != nil {
				return fmt.Errorf("unable to locate project file: %w", err)
			}

			err = srv.run(blend, background)
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

func (srv *Service) run(file *rocketblend.BlendFile, background bool) error {
	cmd, err := srv.driver.GetCMD(file, background, []string{})
	if err != nil {
		return err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to open: %s", err)
	}

	fmt.Println(string(output))

	return nil
}
