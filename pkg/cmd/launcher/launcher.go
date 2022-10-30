package launcher

import (
	"os"

	"github.com/rocketblend/rocketblend/pkg/blendfile"
	"github.com/rocketblend/rocketblend/pkg/client"
)

func Launch() error {
	conf, err := client.LoadConfig()
	if err != nil {
		return err
	}

	client, err := client.NewClient(*conf)
	if err != nil {
		return err
	}

	srv := blendfile.NewService(&blendfile.Config{}, client)

	var path string
	if len((os.Args)) > 1 {
		path = os.Args[1]
	}

	err = srv.Open(path)
	if err != nil {
		return err
	}

	return nil
}
