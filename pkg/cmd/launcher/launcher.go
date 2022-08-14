package launcher

import (
	"fmt"
	"log"
	"os"

	"github.com/rocketblend/rocketblend/pkg/blendfile"
)

func Launch() {
	if len(os.Args) == 1 {
		fmt.Println("No file specified")
		os.Exit(1)
	}

	srv := blendfile.NewService(blendfile.NewConfig(os.Args[2.]))
	blend, err := srv.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if err = srv.Open(blend); err != nil {
		log.Fatal(err)
	}
}
