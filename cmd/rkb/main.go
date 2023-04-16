package main

import (
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend/pkg/rkb"
)

func main() {
	if err := rkb.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
