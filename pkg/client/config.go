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

type NoRocketfileTypeError struct{}

func (m *WrongFileTypeError) Error() string {
	return "File isn't a .blend file"
}

func (m *NoRocketfileTypeError) Error() string {
	return "No .rocketfile file found"
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

	// Validate contents.

	return &rocketfile, nil
}

func FindConfig(path string) (*Rocketfile, error) {
	fileExtension := filepath.Ext(path)
	if fileExtension != ".blend" {
		return nil, &WrongFileTypeError{}
	}

	// Find .rocketfile file in the same directory as the .blend file
	dir := filepath.Dir(path)
	rocketfilePath := filepath.Join(dir, ".rocketfile")

	rocketfile, err := LoadConfig(rocketfilePath)
	if err != nil {
		return nil, err
	}

	return rocketfile, nil
}

func RunConfig(blendFilePath string, rocketfile *Rocketfile) error {
	args := fmt.Sprintf("%s %s", blendFilePath, rocketfile.Args)

	fmt.Println(args)

	cmd := exec.Command("explorer", args)
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
