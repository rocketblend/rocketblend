package client

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/core/executable"
	"github.com/rocketblend/rocketblend/pkg/core/resource"
)

type (
	rocketFile struct {
		Build   string   `json:"build"`
		ARGS    string   `json:"args"`
		Version string   `json:"version"`
		Addons  []string `json:"addons"`
	}

	blendFile struct {
		Exec   *executable.Executable `json:"exec"`
		Path   string                 `json:"path"`
		Addons *[]executable.Addon    `json:"addons"`
		ARGS   string                 `json:"args"`
	}
)

func (c *Client) Open(path string, output string) error {
	file := &blendFile{}
	if path == "" {
		exec, err := c.findDefaultExecutable()
		if err != nil {
			return fmt.Errorf("failed to find default executable: %s", err)
		}

		file.Exec = exec
		file.Addons = &[]executable.Addon{}
	} else {
		loaded, err := c.load(path)
		if err != nil {
			return fmt.Errorf("failed to load blend file: %s", err)
		}

		file = loaded
	}

	switch output {
	case "json":
		json, err := json.Marshal(file)
		if err != nil {
			return fmt.Errorf("failed to marshal blend file: %s", err)
		}
		fmt.Println(string(json))
	case "cmd":
		if err := c.run(file); err != nil {
			return fmt.Errorf("failed to run default build: %s", err)
		}
	default:
		return fmt.Errorf("invalid output format: %s", output)
	}

	return nil
}

func (c *Client) load(path string) (*blendFile, error) {
	ext := filepath.Ext(path)
	if ext != ".blend" {
		return nil, fmt.Errorf("invalid file extension: %s", ext)
	}

	f, err := os.ReadFile(filepath.Join(filepath.Dir(path), "rocketfile.json"))
	if err != nil {
		return nil, fmt.Errorf("failed to read rocketfile: %s", err)
	}

	var rkt rocketFile
	if err := json.Unmarshal(f, &rkt); err != nil {
		return nil, fmt.Errorf("failed to unmarshal rocketfile: %s", err)
	}

	// Get build executable path.
	exec, err := c.findExecutableByBuildReference(rkt.Build)
	if err != nil {
		return nil, fmt.Errorf("failed to find executable: %s", err)
	}

	addons, err := c.getExecutableAddonsByReference(rkt.Addons)
	if err != nil {
		return nil, fmt.Errorf("failed to find all addon directories: %s", err)
	}

	return &blendFile{
		Exec:   exec,
		Path:   path,
		Addons: addons,
		ARGS:   rkt.ARGS,
	}, nil
}

func (c *Client) run(file *blendFile) error {
	addons := append(*file.Exec.Addons, *file.Addons...)
	json, err := json.Marshal(addons)
	if err != nil {
		return fmt.Errorf("failed to marshal addons: %s", err)
	}

	args := []string{}
	if c.conf.Features.Addons {
		script, err := c.FindResource(resource.Startup)
		if err != nil {
			return fmt.Errorf("failed to find startup script: %s", err)
		}

		args = append(args, []string{
			"--python",
			script.OutputPath,
			"--",
			"-a",
			string(json),
		}...)
	}

	if file.Path != "" {
		args = append([]string{file.Path}, args...)
	}

	cmd := exec.Command(file.Exec.Path, args...)

	if c.conf.Debug {
		fmt.Println(strings.ReplaceAll(cmd.String(), "\"", "\\\""))
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open: %s", err)
	}

	return nil
}
