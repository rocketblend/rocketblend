package command

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func (srv *Service) newCreateCommand() *cobra.Command {
	// var build string
	var path string

	c := &cobra.Command{
		Use:   "create",
		Short: "creates a new project",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			path, err := validatePath(path)
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			if err := srv.driver.Create(path); err != nil {
				cmd.PrintErrln(err)
				return
			}
		},
	}

	// c.Flags().StringVarP(&build, "build", "b", "", "build to use to create the file")
	c.Flags().StringVarP(&path, "path", "p", "", "path/name of the file to create")
	if err := c.MarkFlagRequired("path"); err != nil {
		fmt.Println(err)
	}

	return c
}

func validatePath(path string) (string, error) {
	if filepath.Ext(path) != ".blend" {
		return "", fmt.Errorf("path must contain a .blend file")
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absolutePath, nil
}
