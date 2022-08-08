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

func LoadAndOpenConfig(path string) (err error) {
	rocketConfig, err := client.GetBlendConfig(path)
	if err != nil {
		return
	}

	localConfig, err := client.GetLocalConfig()
	if err != nil {
		return
	}

	buildPath, err := client.FindBuildPathFromHash(localConfig, rocketConfig.GetString("build"))
	if err != nil {
		return
	}

	err = client.Open(buildPath, path, rocketConfig.GetString("args"))
	if err != nil {
		return
	}

	return
}
