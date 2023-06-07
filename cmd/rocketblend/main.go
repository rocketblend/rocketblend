package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	cli "github.com/rocketblend/rocketblend/pkg/rocketblend"
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
	signal.Notify(sigs)

	go func() {
		// Wait for a signal
		<-sigs
		// fmt.Println("Cancellation received: ", s)

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
