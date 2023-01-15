package cli

import (
	"github.com/rocketblend/rocketblend/pkg/cmd/cli/command"
	"github.com/rocketblend/rocketblend/pkg/core"
)

func Execute() error {
	driver, err := core.New(nil)
	if err != nil {
		return err
	}

	srv := command.NewService(driver)
	rootCMD := srv.NewCommand()

	return rootCMD.Execute()
}
