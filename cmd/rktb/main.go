package main

import (
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend/pkg/rktb"
)

func main() {
	if err := rktb.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
