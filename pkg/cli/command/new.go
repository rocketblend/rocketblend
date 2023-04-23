package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newNewCommand() *cobra.Command {
	var skipInstall bool

	c := &cobra.Command{
		Use:   "new [name]",
		Short: "Create a new project",
		Long:  `Create a new project`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := srv.validateName(args[0]); err != nil {
				cmd.PrintErrln(err)
				return
			}

			ref := reference.Reference(srv.config.GetValueByString("defaultBuild"))
			err := srv.driver.Create(args[0], srv.flags.workingDirectory, ref, skipInstall)
			if err != nil {
				cmd.Println(err)
				return
			}
		},
	}

	c.Flags().BoolVarP(&skipInstall, "skip-install", "s", false, "skip installing dependencies")

	return c
}

func (srv *Service) validateName(name string) error {
	if filepath.IsAbs(name) || strings.Contains(name, string(filepath.Separator)) {
		return fmt.Errorf("%q is not a valid project name, it should not contain any path separators", name)
	}

	if ext := filepath.Ext(name); ext != "" {
		return fmt.Errorf("%q is not a valid project name, it should not contain any file extension", name)
	}

	return nil
}
