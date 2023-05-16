package main

import (
	"github.com/rocketblend/rocketblend/pkg/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		return
	}
}
