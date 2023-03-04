package command

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/spf13/cobra"
)

func (srv *Service) newRenderCommand() *cobra.Command {
	var background bool

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
			blend, err := srv.findBlendFile()
			if err != nil {
				cmd.Println(err)
				return
			}

			data := map[string]interface{}{
				"PROJECTNAME": "project",
			}

			out, err := srv.parseOutputTemplate(output, data)
			if err != nil {
				cmd.Println(err)
				return
			}

			runArgs := []string{
				fmt.Sprintf("--render-output=%s", out),
				fmt.Sprintf("--frame-start=%d", frameStart),
				fmt.Sprintf("--frame-end=%d", frameEnd),
				fmt.Sprintf("--frame-jump=%d", frameStep),
				fmt.Sprintf("--render-format=%s", format),
				"-x 1", // Set option to add the file extension to the end of the file.
				"-a",   // Render frames from start to end
			}

			err = srv.driver.Run(blend, background, runArgs)
			if err != nil {
				cmd.Println(err)
				return
			}
		},
	}

	c.Flags().BoolVarP(&background, "background", "b", false, "run in the background")

	c.Flags().IntVarP(&frameStart, "frame-start", "s", 0, "start frame")
	c.Flags().IntVarP(&frameEnd, "frame-end", "e", 0, "end frame")
	c.Flags().IntVarP(&frameStep, "frame-step", "t", 1, "frame step")

	c.Flags().StringVarP(&output, "output", "o", "//frames/{{PROJECTNAME}}-######", "set the render path and file name")
	c.Flags().StringVarP(&format, "format", "f", "PNG", "set the render format")

	return c
}

func (srv *Service) parseOutputTemplate(str string, data map[string]interface{}) (string, error) {
	// Define a new template with the input string
	tpl, err := template.New("output").Parse(str)
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
