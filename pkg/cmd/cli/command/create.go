package command

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/spf13/cobra"
)

func (srv *Service) newCreateCommand() *cobra.Command {
	var projectName string
	var projectPath string
	var buildReference string

	var defaultBuild string = srv.driver.GetDefaultBuildReference()

	c := &cobra.Command{
		Use:   "create [flags]",
		Short: "Create a new blender project",
		Long:  `create a new blender project with the specified name, build number and path to create at`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := checkName(projectName); err != nil {
				cmd.PrintErrln(err)
				return
			}

			projectPath, err := checkAndConvertPath(projectPath)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			reference := reference.Reference(buildReference)
			if !reference.IsValid() {
				cmd.PrintErrln(fmt.Errorf("%q is not a valid build reference string", buildReference))
				return
			}

			if err := srv.driver.Create(projectName, projectPath, reference); err != nil {
				cmd.PrintErrln(err)
				return
			}
		},
	}

	c.Flags().StringVarP(&projectName, "name", "n", "", "the name of the project (required)")
	c.MarkFlagRequired("name")

	c.Flags().StringVarP(&buildReference, "build", "b", defaultBuild, "the build reference to use for the project")

	c.Flags().StringVarP(&projectPath, "path", "p", ".", "the path to create the project in")

	return c
}

func checkName(name string) error {
	if filepath.IsAbs(name) || strings.Contains(name, string(filepath.Separator)) {
		return fmt.Errorf("%q is not a valid project name, it should not contain any path separators", name)
	}

	if ext := filepath.Ext(name); ext != "" {
		return fmt.Errorf("%q is not a valid project name, it should not contain any file extension", name)
	}

	return nil
}

func checkAndConvertPath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}

	// get the absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
