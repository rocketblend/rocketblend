package main

import (
	"log"

	"github.com/rocketblend/rocketblend/pkg/cli"
)

func main() {
	app, err := cli.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Execute(); err != nil {
		return
	}
}
