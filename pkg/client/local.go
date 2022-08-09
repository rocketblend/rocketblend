package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/user"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Build struct {
	Platform    string
	Name        string
	Version     string
	Tag         string
	Hash        string
	DownloadUrl string
	CrawledAt   time.Time
}

type Response struct {
	Data []Build `json:"data"`
}

type Install struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Verison string `json:"verison"`
	Hash    string `json:"hash"`
}

type Config struct {
	InstallationDir string    `json:"installationDir"`
	Installed       []Install `json:"installed"`
	Remotes         []string  `json:"remotes"`
}

type buildNotFound struct{}

func (m *buildNotFound) Error() string {
	return "Build not found"
}

const (
	DefaultRemote string = "http://localhost:3000"
	AppName       string = "rocketblend"
	ConfigName    string = "config"
	ConfigType    string = "json"
)

var (
	configViper    = viper.New()
	availableViper = viper.New()

	currentConfig  Config
	availableCache []Build
)

func init() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	configDir := filepath.Join(usr.HomeDir, fmt.Sprintf(".%s", AppName))
	installationDir := filepath.Join(configDir, "installations")

	configViper = viper.New()
	configViper.SetConfigType(ConfigType)
	configViper.SetConfigName(ConfigName)
	configViper.AddConfigPath(configDir)

	configViper.SetDefault("remotes", []string{DefaultRemote})
	configViper.SetDefault("installationDir", installationDir)
	configViper.SetDefault("installed", []Install{})

	if err := configViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, creating one...")
			// Config file not found; ignore error if desired
			if err := configViper.WriteConfig(); err != nil {
				panic("Config file not found and could not be created")
			}
		} else {
			panic(err)
		}
	}

	err = configViper.Unmarshal(&currentConfig)
	if err != nil {
		panic("Unable to decode config into struct")
	}

	fmt.Println(currentConfig)
}

func FetchAvailableBuilds() error {
	response, err := http.Get("http://localhost:3000/builds")

	if err != nil {
		fmt.Print(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	availableCache = responseObject.Data

	return nil
}

func FindBuildPathFromHash(hash string) (buildPath string, err error) {
	fmt.Println(currentConfig.Installed)

	for i := range currentConfig.Installed {
		if currentConfig.Installed[i].Hash == hash {
			buildPath = currentConfig.Installed[i].Path
			break
		}
	}

	if buildPath == "" {
		err = &buildNotFound{}
	}

	return
}
