package launcher

import (
	"fmt"
	"log"
	"os"

	"github.com/rocketblend/rocketblend/pkg/client"
)

func Launch() {
	if len(os.Args) == 1 {
		fmt.Println("No file specified")
		os.Exit(1)
	}

	path := os.Args[1]
	err := LoadAndOpenConfig(path)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadAndOpenConfig(path string) error {
	rocketConfig, err := client.FindConfig(path)
	if err != nil {
		return err
	}

	err = client.RunConfig(path, rocketConfig)
	if err != nil {
		return err
	}

	return nil
}
