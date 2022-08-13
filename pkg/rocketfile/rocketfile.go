package rocketfile

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configType string = "json"
	configName string = "rocketfile"
)

type WrongFileTypeError struct{}

func (m *WrongFileTypeError) Error() string {
	return "File isn't a .blend file"
}

func LoadViperConfig(path string) (rocketfileViper *viper.Viper, err error) {
	rocketfileViper = viper.New()
	rocketfileViper.SetConfigFile(configType)
	rocketfileViper.SetConfigName(configName)
	rocketfileViper.AddConfigPath(path)

	err = rocketfileViper.ReadInConfig()
	if err != nil {
		return
	}

	return
}

func GetBlendConfig(path string) (rocketfileViper *viper.Viper, err error) {
	fileExtension := filepath.Ext(path)
	if fileExtension != ".blend" {
		return nil, &WrongFileTypeError{}
	}

	// Move to init.
	// Find .rocketfile file in the same directory as the .blend file
	dir := filepath.Dir(path)
	rocketfileViper, err = LoadViperConfig(dir)
	if err != nil {
		return
	}

	return
}

func Open(buildPath string, blendPath string, args string) (err error) {
	combinedArgs := fmt.Sprintf("%s %s", blendPath, args)

	cmd := exec.Command(buildPath, combinedArgs)
	err = cmd.Start()

	if err != nil {
		return err
	}

	return nil
}

func init() {
	println("Loading rocketfile config")
}
