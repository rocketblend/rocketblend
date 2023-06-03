package main

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/cli"
)

func main() {
	app, err := cli.New()
	if err != nil {
		fmt.Println("Error creating cli app: ", err)
		return
	}

	if err := app.Execute(); err != nil {
		return
	}
}
