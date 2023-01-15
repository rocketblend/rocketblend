package command

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

func (srv *Service) newOpenCommand() *cobra.Command {
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

			if err := srv.open(path, output); err != nil {
				cmd.PrintErrln(err)
			}
		},
	}

	c.Flags().StringVarP(&path, "path", "p", "", "The file path to a .blend file.")
	c.Flags().StringVarP(&output, "output", "o", "cmd", "Output type of the command")
	c.Flags().BoolVarP(&auto, "auto", "a", false, "Enables or disables the automatic detection of .blend files in the current directory.")

	return c
}

func (srv *Service) open(path string, output string) error {
	file, err := srv.driver.Load(path)
	if err != nil {
		return err
	}

	switch output {
	case "json":
		json, err := json.Marshal(file)
		if err != nil {
			return fmt.Errorf("failed to marshal blend file: %s", err)
		}
		fmt.Println(string(json))
	case "cmd":
		if err := srv.driver.Run(file); err != nil {
			return fmt.Errorf("failed to run default build: %s", err)
		}
	default:
		return fmt.Errorf("invalid output format: %s", output)
	}

	return nil
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
