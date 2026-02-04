package command

import (
	"github.com/rocketblend/rocketblend/internal/worker"
	"github.com/spf13/cobra"
)

type WorkerOptions struct {
	commandOpts
	ServerURL  string
	Hostname   string
	WorkerPort int // Local TCP port for generic connector
}

func newWorkerCommand(commandOpts commandOpts) *cobra.Command {
	opts := &WorkerOptions{
		commandOpts: commandOpts,
	}

	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Start a worker",
		Long:  `Starts a worker agent that connects to a host and executes jobs by managing local Blender instances.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			l := getLogger(opts.Global.Level, opts.Global.Verbose)

			l.Info("starting worker", map[string]any{
				"server": opts.ServerURL,
			})

			workerOpts := []worker.Option{
				worker.WithServerURL(opts.ServerURL),
				worker.WithWorkingDir(opts.Global.WorkingDirectory),
				worker.WithLogger(l),
			}

			if opts.Hostname != "" {
				workerOpts = append(workerOpts, worker.WithHostname(opts.Hostname))
			}

			w := worker.New(workerOpts...)

			return w.Start(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&opts.ServerURL, "url", "http://localhost:8080", "URL of the Server")
	cmd.Flags().StringVar(&opts.Hostname, "hostname", "", "Hostname to identify the worker (defaults to OS hostname)")

	return cmd
}
