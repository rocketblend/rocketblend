package core

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/core/resource"
	"github.com/rocketblend/rocketblend/pkg/core/rocketfile"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/semver"
)

type (
	Addon struct {
		Name    string         `json:"name"`
		Version semver.Version `json:"version"`
		Path    string         `json:"path"`
	}

	Build struct {
		Path   string   `json:"path"`
		Addons *[]Addon `json:"addons"`
		ARGS   string   `json:"args"`
	}

	BlendFile struct {
		Build  *Build   `json:"build"`
		Path   string   `json:"path"`
		Addons *[]Addon `json:"addons"`
		ARGS   string   `json:"args"`
	}
)

func (d *Driver) Load(path string) (*BlendFile, error) {
	file := &BlendFile{}
	if path == "" {
		ref, err := d.getDefaultBuild()
		if err != nil {
			return nil, err
		}

		build, err := d.findBuildByReference(ref)
		if err != nil {
			return nil, err
		}

		file.Build = build
		file.Addons = &[]Addon{}
	} else {
		loaded, err := d.load(path)
		if err != nil {
			return nil, err
		}

		file = loaded
	}

	return file, nil
}

func (d *Driver) Run(file *BlendFile) error {
	args := []string{}
	if d.conf.Features.Addons {
		addons := append(*file.Build.Addons, *file.Addons...)
		json, err := json.Marshal(addons)
		if err != nil {
			return fmt.Errorf("failed to marshal addons: %s", err)
		}

		script, err := d.resource.FindByName(resource.Startup)
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

	cmd := exec.Command(file.Build.Path, args...)

	if d.conf.Debug {
		fmt.Println(strings.ReplaceAll(cmd.String(), "\"", "\\\""))
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open: %s", err)
	}

	return nil
}

func (d *Driver) getDefaultBuild() (string, error) {
	build := d.conf.Defaults.Build
	if build == "" {
		return "", fmt.Errorf("no default build set")
	}

	return build, nil
}

func (d *Driver) load(path string) (*BlendFile, error) {
	ext := filepath.Ext(path)
	if ext != ".blend" {
		return nil, fmt.Errorf("invalid file extension: %s", ext)
	}

	rkt, err := rocketfile.Load(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

	// Get build executable path.
	build, err := d.findBuildByReference(rkt.Build)
	if err != nil {
		return nil, fmt.Errorf("failed to find executable: %s", err)
	}

	addons, err := d.getAddonsByReference(rkt.Addons)
	if err != nil {
		return nil, fmt.Errorf("failed to find all addon directories: %s", err)
	}

	return &BlendFile{
		Build:  build,
		Path:   path,
		Addons: addons,
		ARGS:   rkt.ARGS,
	}, nil
}

func (d *Driver) findBuildByReference(ref string) (*Build, error) {
	pack, err := d.pack.FindByReference(reference.Reference(ref))
	if err != nil {
		return nil, fmt.Errorf("failed to find build: %s", err)
	}

	if pack.Build == nil {
		return nil, fmt.Errorf("packge has no build")
	}

	addons, err := d.getAddonsByReference(pack.Build.Addons)
	if err != nil {
		return nil, fmt.Errorf("failed to find all addons for build: %s", err)
	}

	return &Build{
		Path:   filepath.Join(d.conf.Directories.Installations, ref, pack.Build.GetSourceForPlatform(d.conf.Platform).Executable),
		Addons: addons,
		ARGS:   pack.Build.Args,
	}, nil
}

func (d *Driver) getAddonsByReference(ref []string) (*[]Addon, error) {
	addons := []Addon{}
	if d.conf.Features.Addons {
		for _, r := range ref {
			addon, err := d.getAddonByReference(r)
			if err != nil {
				return nil, fmt.Errorf("failed to find addon: %s", err)
			}

			addons = append(addons, *addon)
		}
	}

	return &addons, nil
}

func (d *Driver) getAddonByReference(ref string) (*Addon, error) {
	pack, err := d.pack.FindByReference(reference.Reference(ref))
	if err != nil {
		return nil, fmt.Errorf("failed to find addon: %s", err)
	}

	if pack.Addon == nil {
		return nil, fmt.Errorf("packge has no addon")
	}

	return &Addon{
		Name:    pack.Addon.Name,
		Version: pack.Addon.Version,
		Path:    filepath.Join(d.conf.Directories.Installations, ref, pack.Addon.Source.File),
	}, nil
}
