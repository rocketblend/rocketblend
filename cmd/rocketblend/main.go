package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rocketblend/rocketblend/internal/cli"
)

func main() {
	app, err := cli.New()
	if err != nil {
		fmt.Println("Error creating cli app: ", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal catching
	sigs := make(chan os.Signal, 1)

	// Catch all signals since we can't block SIGKILL
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Wait for a signal
		<-sigs

		// Cancel the context on receipt of a signal
		cancel()
	}()

	if err := app.ExecuteContext(ctx); err != nil {
		if ctx.Err() == context.Canceled {
			return
		}

		fmt.Println("Error executing: ", err)
	}
}
