package command

import (
	"github.com/rocketblend/rocketblend/internal/orchestration/server"
	"github.com/spf13/cobra"
)

type ServeOptions struct {
	commandOpts
	Port int
	Host string
}

func newServeCommand(commandOpts commandOpts) *cobra.Command {
	opts := &ServeOptions{
		commandOpts: commandOpts,
	}

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the orchestrator host",
		Long:  `Starts the central HTTP/WebSocket server that manages the job queue and worker pool.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			l := getLogger(opts.Global.Level, opts.Global.Verbose)
			l.Info("starting host", map[string]any{
				"host": opts.Host,
				"port": opts.Port,
			})

			srv := server.New(
				server.WithPort(opts.Port),
				server.WithHost(opts.Host),
				server.WithLogger(l),
			)

			return srv.Start(cmd.Context())
		},
	}

	cmd.Flags().IntVarP(&opts.Port, "port", "p", 8080, "Port to listen on")
	cmd.Flags().StringVarP(&opts.Host, "host", "H", "0.0.0.0", "Host interface to listen on")

	return cmd
}
