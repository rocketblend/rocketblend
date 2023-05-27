package command

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rocketblend/rocketblend/pkg/blenderparser"
	"github.com/rocketblend/rocketblend/pkg/rocketblend"
	"github.com/spf13/cobra"
)

type templateData struct {
	Project string
}

// newRenderCommand creates a new cobra command for rendering the project.
// It sets up all necessary flags and executes the rendering through the driver.
func (srv *Service) newRenderCommand() *cobra.Command {
	var frameStart int
	var frameEnd int
	var frameStep int

	var output string
	var format string

	c := &cobra.Command{
		Use:   "render",
		Short: "Renders the project",
		Long:  `Renders the project from the specified start frame to the end frame, with the given step. Outputs the render in the provided format.`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if frameEnd < frameStart || frameStep <= 0 {
				return fmt.Errorf("invalid frame range or step")
			}

			blend, err := srv.findBlendFile(srv.flags.workingDirectory)
			if err != nil {
				return fmt.Errorf("failed to find blend file: %w", err)
			}

			name := filepath.Base(blend.Path)
			data := templateData{
				Project: strings.TrimSuffix(name, filepath.Ext(name)),
			}

			out, err := srv.parseOutputTemplate(output, data)
			if err != nil {
				return fmt.Errorf("failed to parse output template: %w", err)
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

			err = srv.run(blend, true, runArgs)
			if err != nil {
				return fmt.Errorf("failed to run driver: %w", err)
			}

			return nil
		},
	}

	c.Flags().IntVarP(&frameStart, "frame-start", "s", 0, "start frame")
	c.Flags().IntVarP(&frameEnd, "frame-end", "e", 0, "end frame")
	c.Flags().IntVarP(&frameStep, "frame-step", "t", 1, "frame step")

	c.Flags().StringVarP(&output, "output", "o", "//output/{{.Project}}-#####", "set the render path and file name")
	c.Flags().StringVarP(&format, "format", "f", "PNG", "set the render format")

	return c
}

func (srv *Service) run(file *rocketblend.BlendFile, background bool, args []string) error {
	cmd, err := srv.driver.GetCMD(file, background, args)
	if err != nil {
		return err
	}

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("creating stdout pipe: %w", err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			info, err := blenderparser.ParseRenderOutput(scanner.Text())
			if err != nil {
				fmt.Println("Error parsing blender output:", err)
			} else {
				// Print or use the render info
				fmt.Printf("Frame: %d, Memory: %s, Peak Memory: %s, Time: %s, Operation: %s\n",
					info.FrameNumber, info.Memory, info.PeakMemory, info.Time, info.Operation)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("waiting for command: %w", err)
	}

	return nil
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

	return buf.String(), nil
}
