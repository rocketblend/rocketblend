package main

import (
	"fmt"

	cli "github.com/rocketblend/rocketblend/pkg/rocketblend"
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
