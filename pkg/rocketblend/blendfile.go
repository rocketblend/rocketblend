package rocketblend

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketfile"
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

func (d *Driver) Create(ctx context.Context, name string, path string, reference reference.Reference, skipDeps bool) error {
	rkt := rocketfile.RocketFile{
		Build: reference.String(),
	}

	if err := rocketfile.Save(path, &rkt); err != nil {
		return fmt.Errorf("failed to create rocketfile: %s", err)
	}

	if !skipDeps {
		err := d.InstallDependencies(ctx, path, nil, false)
		if err != nil {
			return err
		}
	}

	// TODO: convert all functions to use reference.Reference
	build, err := d.findBuildByReference(reference.String())
	if err != nil {
		return err
	}

	blendFile := &BlendFile{
		Build: build,
		Path:  filepath.Join(path, name+BlenderFileExtension),
	}

	if err := d.create(blendFile); err != nil {
		return err
	}

	return nil
}

func (d *Driver) Load(path string) (*BlendFile, error) {
	file := &BlendFile{}
	if path == "" {
		build, err := d.getDefaultBuild()
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

func (d *Driver) GetCMD(ctx context.Context, file *BlendFile, background bool, postArgs []string) (*exec.Cmd, error) {
	preArgs := []string{}
	if background {
		preArgs = append(preArgs, "-b")
	}

	if file.Path != "" {
		preArgs = append(preArgs, []string{file.Path}...)
	}

	if d.addonsEnabled {
		addons := append(*file.Build.Addons, *file.Addons...)
		json, err := json.Marshal(addons)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal addons: %s", err)
		}

		postArgs = append([]string{
			"--python-expr",
			d.resource.GetAddonScript(),
		}, postArgs...)

		postArgs = append(postArgs, []string{
			"--",
			"-a",
			string(json),
		}...)
	}

	// Blender requires arguments to be in a specific order
	args := append(preArgs, postArgs...)
	cmd := exec.CommandContext(ctx, file.Build.Path, args...)

	if d.debug {
		fmt.Println(strings.ReplaceAll(cmd.String(), "\"", "\\\""))
	}

	return cmd, nil
}

func (d *Driver) getDefaultBuild() (*Build, error) {
	ref := d.defaultBuild
	if ref == "" {
		return nil, fmt.Errorf("no default build set")
	}

	build, err := d.findBuildByReference(ref)
	if err != nil {
		return nil, err
	}

	return build, nil
}

func (d *Driver) load(path string) (*BlendFile, error) {
	ext := filepath.Ext(path)
	if ext != BlenderFileExtension {
		return nil, fmt.Errorf("invalid file extension: %s", ext)
	}

	rkt, err := rocketfile.Load(filepath.Dir(path))
	if err != nil {
		return nil, err
	}

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
		Path:   filepath.Join(d.InstallationDirectory, ref, pack.Build.GetSourceForPlatform(d.platform).Executable),
		Addons: addons,
		ARGS:   pack.Build.Args,
	}, nil
}

func (d *Driver) getAddonsByReference(ref []string) (*[]Addon, error) {
	addons := []Addon{}
	if d.addonsEnabled {
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

	var path string
	if pack.Addon.Source != nil {
		path = filepath.Join(d.InstallationDirectory, ref, pack.Addon.Source.File)
	}

	return &Addon{
		Name:    pack.Addon.Name,
		Version: *pack.Addon.Version,
		Path:    path,
	}, nil
}

func (d *Driver) create(blendFile *BlendFile) error {
	script, err := d.resource.GetCreateScript(blendFile.Path)
	if err != nil {
		return err
	}

	cmd := exec.Command(blendFile.Build.Path, "-b", "--python-expr", script)

	if d.debug {
		fmt.Println(strings.ReplaceAll(cmd.String(), "\"", "\\\""))
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create: %s", err)
	}

	return nil
}
