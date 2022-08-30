package launcher

import (
	"fmt"
	"os"

	"github.com/rocketblend/rocketblend/pkg/blendfile"
	"github.com/rocketblend/rocketblend/pkg/client"
)

func Launch() error {
	if len(os.Args) == 1 {
		return fmt.Errorf("no file specified")
	}

	conf, err := client.LoadConfig()
	if err != nil {
		return err
	}

	client, err := client.NewClient(*conf)
	if err != nil {
		return err
	}

	srv := blendfile.NewService(&blendfile.Config{}, client)
	blend, err := srv.Load(os.Args[1])
	if err != nil {
		return err
	}

	if err = srv.Open(blend); err != nil {
		return err
	}

	return nil
}
