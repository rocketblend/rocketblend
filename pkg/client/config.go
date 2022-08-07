package client

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Rocketfile struct {
	Version   string `json:"version"`
	BlendFile string `json:"blendfile"`
	Build     string `json:"build"`
	Name      string `json:"name"`
	Args      string `json:"args"`
}

type WrongFileTypeError struct{}

func (m *WrongFileTypeError) Error() string {
	return "File isn't a .rocketfile"
}

func LoadConfig(path string) (*Rocketfile, error) {
	fileExtension := filepath.Ext(path)
	if fileExtension != ".rocketfile" {
		return nil, &WrongFileTypeError{}
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var rocketfile Rocketfile
	json.Unmarshal([]byte(byteValue), &rocketfile)

	// Convert relative path to absolute path
	rocketfile.BlendFile = filepath.Join(filepath.Dir(path), rocketfile.BlendFile)

	// Validate contents.
	fmt.Println(rocketfile)

	return &rocketfile, nil
}

func RunConfig(rocketfile *Rocketfile) error {
	cmd := exec.Command("explorer", rocketfile.BlendFile)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
