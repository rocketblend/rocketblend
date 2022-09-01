package client

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rocketblend/rocketblend/pkg/core/addon"
	"github.com/rocketblend/rocketblend/pkg/core/executable"
	"github.com/rocketblend/rocketblend/pkg/core/install"
	"github.com/rocketblend/rocketblend/pkg/core/library"
	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/scribble"
)

type (
	InstallService interface {
		FindAll() ([]*install.Install, error)
		FindByID(id string) (*install.Install, error)
		Create(i *install.Install) error
		Remove(id string) error
	}

	AddonService interface {
		FindAll() ([]*addon.Addon, error)
		FindByID(id string) (*addon.Addon, error)
		Create(i *addon.Addon) error
		Remove(id string) error
	}

	LibraryService interface {
		FindBuildByPath(path string) (*library.Build, error)
		FindPackageByPath(path string) (*library.Package, error)
		FetchBuild(str string) (*library.Build, error)
		FetchPackage(str string) (*library.Package, error)
	}

	DownloadService interface {
		Download(url string, path string) error
	}

	ArchiverService interface {
		Extract(path string) error
	}

	EncoderService interface {
		Hash(str string) string
	}

	Config struct {
		DBDir           string
		InstallationDir string
		Platform        runtime.Platform
	}

	Client struct {
		install    InstallService
		addon      AddonService
		library    LibraryService
		downloader DownloadService
		archiver   ArchiverService
		encoder    EncoderService
		conf       Config
	}
)

func NewClient(conf Config) (*Client, error) {
	db, err := scribble.New(conf.DBDir, nil)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(conf.InstallationDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create installation directory: %w", err)
	}

	client := &Client{
		install:    NewInstallService(db),
		addon:      NewAddonService(db),
		library:    NewLibraryService(),
		downloader: NewDownloaderService(conf.InstallationDir),
		archiver:   NewArchiverService(true),
		encoder:    NewEncoderService(),
		conf:       conf,
	}

	return client, nil
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find user home directory: %v", err)
	}

	platform := runtime.DetectPlatform()
	if platform == runtime.Undefined {
		return nil, fmt.Errorf("cannot detect platform")
	}

	appDir := filepath.Join(home, fmt.Sprintf(".%s", "rocketblend"))
	conf := Config{
		InstallationDir: filepath.Join(appDir, "installations"),
		DBDir:           filepath.Join(appDir, "data"),
		Platform:        platform,
	}

	return &conf, nil
}

func (c *Client) Platform() runtime.Platform {
	return c.conf.Platform
}

func (c *Client) FindInstall(ref string) (*install.Install, error) {
	id := c.encoder.Hash(ref)
	install, err := c.install.FindByID(id)
	if err != nil {
		return nil, err
	}

	return install, nil
}

func (c *Client) FindAddon(ref string) (*addon.Addon, error) {
	id := c.encoder.Hash(ref)
	addon, err := c.addon.FindByID(id)
	if err != nil {
		return nil, err
	}

	return addon, nil
}

func (c *Client) AddInstall(path string) error {
	build, err := c.library.FindBuildByPath(path)
	if err != nil {
		return err
	}

	i, _ := c.FindInstall(build.Reference)
	if i != nil {
		return fmt.Errorf("build already installed")
	}

	err = c.install.Create(c.newInstall(build.Reference, path, build.Packages))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AddAddon(path string) error {
	pack, err := c.library.FindPackageByPath(path)
	if err != nil {
		return err
	}

	a, _ := c.FindAddon(pack.Reference)
	if a != nil {
		return fmt.Errorf("addon already installed")
	}

	err = c.addon.Create(c.newAddon(pack.Reference, path))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) InstallBuild(ref string) error {
	// Check if install already exists
	inst, _ := c.FindInstall(ref)
	if inst != nil {
		return fmt.Errorf("already installed")
	}

	// Fetch build from ref
	build, err := c.library.FetchBuild(ref)
	if err != nil {
		return err
	}

	if build == nil {
		return fmt.Errorf("invalid build")
	}

	// Output directory
	outPath := filepath.Join(c.conf.InstallationDir, ref)

	// Create output directories
	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		return err
	}

	// build info for current platform
	source := build.GetSourceForPlatform(c.conf.Platform)
	if source == nil {
		return fmt.Errorf("no source found for platform %s", c.conf.Platform)
	}

	// Download URL
	downloadURL := source.URL

	// Download file path
	name := filepath.Base(downloadURL)
	filePath := filepath.Join(outPath, name)

	// Download file to file path
	err = c.downloader.Download(downloadURL, filePath)
	if err != nil {
		return err
	}

	// Extract the archived file
	if err := c.archiver.Extract(filePath); err != nil {
		return err
	}

	// Markshal build
	data, err := json.Marshal(build)
	if err != nil {
		return err
	}

	// Write out build.json
	if err := os.WriteFile(filepath.Join(outPath, library.BuildFile), data, os.ModePerm); err != nil {
		return err
	}

	// Add install to database
	err = c.install.Create(c.newInstall(ref, outPath, build.Packages))
	if err != nil {
		return err
	}

	// TODO: call asynchronously
	// Install packages
	for _, p := range build.Packages {
		err = c.installPackageIgnorable(p, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) InstallPackage(ref string) error {
	return c.installPackageIgnorable(ref, false)
}

func (c *Client) RemoveInstall(hash string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) RemoveAddon(hash string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) FindAllInstalls() ([]*install.Install, error) {
	ints, err := c.install.FindAll()
	if err != nil {
		return nil, err
	}

	return ints, nil
}

func (c *Client) FindAllAddons() ([]*addon.Addon, error) {
	addons, err := c.addon.FindAll()
	if err != nil {
		return nil, err
	}

	return addons, nil
}

func (c *Client) FindBuildByPath(build string) (*library.Build, error) {
	return c.library.FindBuildByPath(build)
}

func (c *Client) FetchPackageByPath(pack string) (*library.Package, error) {
	return c.library.FindPackageByPath(pack)
}

func (c *Client) OpenProject(file string, ref string, args string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) CreateProject(name string, path string, ref string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) FindExecutableByBuildReference(ref string) (*executable.Executable, error) {
	install, err := c.FindInstall(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to find install: %s", err)
	}

	build, err := c.library.FindBuildByPath(install.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to find build: %s", err)
	}

	addons, err := c.FindAllAddonDirectories(build.Packages)
	if err != nil {
		return nil, fmt.Errorf("failed to find all addons for build: %s", err)
	}

	return &executable.Executable{
		Path:   filepath.Join(install.Path, build.GetSourceForPlatform(c.conf.Platform).Executable),
		Addons: addons,
		ARGS:   build.Args,
	}, nil
}

func (c *Client) FindAllAddonDirectories(ref []string) ([]string, error) {
	var addons []string
	for _, a := range ref {
		dir, err := c.findAddonDirectoryByPackageReference(a)
		if err != nil {
			return nil, fmt.Errorf("failed to find addon directory: %s", err)
		}

		addons = append(addons, dir)
	}

	return addons, nil
}

func (c *Client) findAddonDirectoryByPackageReference(ref string) (string, error) {
	addon, err := c.FindAddon(ref)
	if err != nil {
		return "", fmt.Errorf("failed to find addon: %s", err)
	}

	pack, err := c.library.FindPackageByPath(addon.Path)
	if err != nil {
		return "", fmt.Errorf("failed to find package: %s", err)
	}

	return filepath.Join(addon.Path, pack.Source.File), nil
}

func (c *Client) installPackageIgnorable(ref string, ignore bool) error {
	// TODO: Move downloading packages/builds into library service.

	// Check if addon already exists
	adn, _ := c.FindAddon(ref)
	if adn != nil {
		if !ignore {
			return fmt.Errorf("already installed")
		}
		return nil
	}

	// Fetch package from ref
	pack, err := c.library.FetchPackage(ref)
	if err != nil {
		return err
	}

	if pack == nil {
		return fmt.Errorf("invalid package")
	}

	// Output directory
	outPath := filepath.Join(c.conf.InstallationDir, ref)

	// Create output directories
	err = os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		return err
	}

	// download file path
	name := filepath.Base(pack.Source.URL)
	filePath := filepath.Join(outPath, name)

	// Download file to file path
	err = c.downloader.Download(pack.Source.URL, filePath)
	if err != nil {
		return err
	}

	// Markshal pack
	data, err := json.Marshal(pack)
	if err != nil {
		return err
	}

	// Write out package.json
	if err := os.WriteFile(filepath.Join(outPath, library.PackgeFile), data, os.ModePerm); err != nil {
		return err
	}

	// Add addon to database
	err = c.addon.Create(c.newAddon(ref, outPath))
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) newInstall(ref string, path string, packs []string) *install.Install {
	return &install.Install{
		Id:       c.encoder.Hash(ref),
		Build:    ref,
		Path:     path,
		Packages: packs,
	}
}

func (c *Client) newAddon(ref string, path string) *addon.Addon {
	return &addon.Addon{
		Id:      c.encoder.Hash(ref),
		Package: ref,
		Path:    path,
	}
}
