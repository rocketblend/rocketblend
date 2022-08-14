package main

import (
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend/pkg/cmd/launcher"
)

func main() {
	if err := launcher.Launch(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
