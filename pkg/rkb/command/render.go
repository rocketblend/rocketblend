package command

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type templateData struct {
	Project string
}

func (srv *Service) newRenderCommand() *cobra.Command {
	var frameStart int
	var frameEnd int
	var frameStep int

	var output string
	var format string

	c := &cobra.Command{
		Use:   "render",
		Short: "Render project",
		Long:  `Render project`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			blend, err := srv.findBlendFile(srv.flags.workingDirectory)
			if err != nil {
				cmd.Println(err)
				return
			}

			name := filepath.Base(blend.Path)
			data := templateData{
				Project: strings.TrimSuffix(name, filepath.Ext(name)),
			}

			out, err := srv.parseOutputTemplate(output, data)
			if err != nil {
				cmd.Println(err)
				return
			}

			runArgs := []string{
				"--frame-start",
				fmt.Sprint(frameStart),
				"--frame-end",
				fmt.Sprint(frameEnd),
				"--frame-jump",
				fmt.Sprint(frameStep),
				"--render-output",
				out,
				"--render-format",
				format,
				"-x", // Set option to add the file extension to the end of the file.
				"1",
				"-a", // Render frames from start to end
			}

			err = srv.driver.Run(blend, true, runArgs)
			if err != nil {
				cmd.Println(err)
				return
			}
		},
	}

	c.Flags().IntVarP(&frameStart, "frame-start", "s", 0, "start frame")
	c.Flags().IntVarP(&frameEnd, "frame-end", "e", 0, "end frame")
	c.Flags().IntVarP(&frameStep, "frame-step", "t", 1, "frame step")

	c.Flags().StringVarP(&output, "output", "o", "//output/{{.Project}}-#####", "set the render path and file name")
	c.Flags().StringVarP(&format, "format", "f", "PNG", "set the render format")

	return c
}

func (srv *Service) parseOutputTemplate(str string, data interface{}) (string, error) {
	// Define a new template with the input string
	tpl, err := template.New("").Parse(str)
	if err != nil {
		return "", err
	}

	// Execute the template with the data object and capture the output
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	// Return the output string
	return buf.String(), nil
}
