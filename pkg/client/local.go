package client

import (
	"github.com/spf13/viper"
)

type install struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Verison string `json:"verison"`
	Hash    string `json:"hash"`
}

type buildNotFound struct{}

func (m *buildNotFound) Error() string {
	return "Build not found"
}

func GetLocalConfig() (config *viper.Viper, err error) {
	config = viper.New()
	config.SetConfigType("json")
	config.SetConfigName("config")
	config.AddConfigPath("C:\\Users\\Elliot\\.rocketblend")

	err = config.ReadInConfig()
	if err != nil {
		return
	}

	return
}

func FindBuildPathFromHash(config *viper.Viper, hash string) (buildPath string, err error) {
	// Move out. Just take type as argument.
	var installs []install
	err = config.UnmarshalKey("installed", &installs)
	if err != nil {
		return
	}

	for i := range installs {
		if installs[i].Hash == hash {
			buildPath = installs[i].Path
			break
		}
	}

	if buildPath == "" {
		err = &buildNotFound{}
	}

	return
}
