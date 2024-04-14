package blender

import (
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/python"
)

type (
	createScriptData struct {
		Path string `json:"path"`
	}
)

func getCreateScript(data *createScriptData) (string, error) {
	result, err := helpers.ParseTemplateWithData(python.CreateScript, data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func getStartupScript() string {
	return python.StartupScript
}
