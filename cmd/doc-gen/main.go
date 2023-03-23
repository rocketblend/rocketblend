package main

import (
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend/pkg/cmd/cli"
)

func main() {
	if err := cli.GenerateDocs("./docs/reference/cli/commands"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
