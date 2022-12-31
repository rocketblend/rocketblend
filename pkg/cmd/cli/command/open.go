package command

import (
	"fmt"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/client"
	"github.com/spf13/cobra"
)

func NewOpenCommand(srv *client.Client) *cobra.Command {
	var path string
	var output string
	var auto bool

	c := &cobra.Command{
		Use:   "open",
		Short: "Opens blender with the specified version",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			if auto && path == "" {
				file, err := findBlendFile()
				if err != nil {
					cmd.PrintErrln(err)
					return
				}

				path = file
			}

			if err := srv.Open(path, output); err != nil {
				cmd.PrintErrln(err)
			}
		},
	}

	c.Flags().StringVarP(&path, "path", "p", "", "The file path to a .blend file.")
	c.Flags().StringVarP(&output, "output", "o", "cmd", "Output type of the command")
	c.Flags().BoolVarP(&auto, "auto", "a", false, "Enables or disables the automatic detection of .blend files in the current directory.")

	return c
}

func findBlendFile() (string, error) {
	// Get a list of all files in the current directory.
	files, err := filepath.Glob("*")
	if err != nil {
		return "", fmt.Errorf("failed to list files in current directory: %w", err)
	}

	// Iterate through the list of files and check if any have a .blend extension.
	for _, file := range files {
		if filepath.Ext(file) == ".blend" {
			// Found a .blend file. Return the full path.
			return filepath.Abs(file)
		}
	}

	// No .blend files found. Return an error.
	return "", fmt.Errorf("no .blend files found in current directory")
}
