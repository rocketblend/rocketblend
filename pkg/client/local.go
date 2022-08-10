package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/user"
	"path/filepath"
	"sort"
	"time"

	"github.com/blang/semver/v4"
	"github.com/spf13/viper"
	"go.lsp.dev/uri"
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
	Version string `json:"version"`
	Hash    string `json:"hash"`
}

type Available struct {
	Hash    string
	Name    string
	Uri     uri.URI
	Version semver.Version
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
	DefaultRemote string = "http://127.0.0.1:3000/builds/"
	AppName       string = "rocketblend"
	ConfigName    string = "config"
	ConfigType    string = "json"
)

var (
	configViper = viper.New()
	//currentConfig Config
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

	// err = configViper.Unmarshal(&currentConfig)
	// if err != nil {
	// 	panic("Unable to decode config into struct")
	// }

	// fmt.Println(currentConfig)
}

func GetInstallationDir() string {
	return configViper.GetString("installationDir")
}

func GetRemotes() []string {
	return configViper.GetStringSlice("remotes")
}

func GetInstalledBuilds() []Install {
	installs := []Install{}
	configViper.UnmarshalKey("installed", &installs)
	return installs
}

func FetchAvailableBuildsFromRemote(url string) ([]Build, error) {
	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Data, nil
}

func FetchAvailableBuildsFromRemotes() ([]Build, error) {
	remotes := configViper.GetStringSlice("remotes")
	availableBuilds := []Build{}

	for _, remote := range remotes {
		builds, err := FetchAvailableBuildsFromRemote(remote)
		if err != nil {
			return nil, err
		}

		availableBuilds = append(availableBuilds, builds...)
	}

	return availableBuilds, nil
}

func GetAvilableBuilds() []Available {
	builds, err := FetchAvailableBuildsFromRemotes()
	if err != nil {
		log.Fatal(err)
	}

	available := []Available{}
	installed := GetInstalledBuilds()
	for _, installed := range installed {
		temp := Available{}
		temp.Hash = installed.Hash
		temp.Name = installed.Name
		temp.Version, _ = semver.Parse(installed.Version)
		temp.Uri, _ = uri.Parse(installed.Path)
		available = append(available, temp)
	}

	for _, build := range builds {
		isExisting := false
		for _, existing := range available {
			if build.Hash == existing.Hash {
				isExisting = true
				break
			}
		}

		if !isExisting {
			temp := Available{}
			temp.Hash = build.Hash
			temp.Name = build.Name
			temp.Version, _ = semver.Parse(build.Version)
			temp.Uri, _ = uri.Parse(build.DownloadUrl)
			available = append(available, temp)
		}
	}

	sort.SliceStable(available, func(i, j int) bool {
		return available[i].Version.GT(available[j].Version)
	})

	return available
}

func FindAvailableBuildFromHash(hash string) *Available {
	available := GetAvilableBuilds()
	for _, build := range available {
		if build.Hash == hash {
			return &build
		}
	}

	return nil
}

func FindInstalledBuildPathFromHash(hash string) (string, error) {
	for _, build := range GetInstalledBuilds() {
		if build.Hash == hash {
			return build.Path, nil
		}
	}

	return "", &buildNotFound{}
}

func AddInstall(install Install) {
	installed := GetInstalledBuilds()
	installed = append(installed, install)
	configViper.Set("installed", installed)
	configViper.WriteConfig()
}
