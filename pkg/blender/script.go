package blender

import (
	"github.com/rocketblend/rocketblend/pkg/helpers"
	"github.com/rocketblend/rocketblend/pkg/python"
)

type (
	createBlendFileData struct {
		FilePath string `json:"filePath"`
	}
)

func createBlendFileScript(data *createBlendFileData) (string, error) {
	result, err := helpers.ParseTemplateWithData(python.CreateScript, data)
	if err != nil {
		return "", err
	}

	return result, nil
}

func startupScript() string {
	return python.StartupScript
}
