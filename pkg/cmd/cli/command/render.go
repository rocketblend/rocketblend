package command

import "github.com/spf13/cobra"

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
			cmd.Println("render called")
		},
	}

	c.Flags().IntVarP(&frameStart, "frame-start", "s", 0, "start frame")
	c.Flags().IntVarP(&frameEnd, "frame-end", "e", 0, "end frame")
	c.Flags().IntVarP(&frameStep, "frame-step", "t", 0, "frame step")

	c.Flags().StringVarP(&output, "output", "o", "", "output file name")
	c.Flags().StringVarP(&format, "format", "f", "", "output format")

	return c
}
