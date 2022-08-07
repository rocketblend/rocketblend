package launcher

import (
	"log"
	"os"

	"github.com/rocketblend/rocketblend/pkg/client"
)

func Launch() {
	var path string = "./.rocketfile"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	err := LoadAndOpenConfig(path)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadAndOpenConfig(path string) error {
	rocketConfig, err := client.LoadConfig(path)
	if err != nil {
		return err
	}

	err = client.RunConfig(rocketConfig)
	if err != nil {
		return err
	}

	return nil
}
