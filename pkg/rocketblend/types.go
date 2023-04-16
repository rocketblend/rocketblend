package rocketblend

import (
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
)

const (
	Name                 = "rocketblend"
	BlenderFileExtension = ".blend"
)

type (
	// Logger is a generic logger interface
	Logger interface {
		Debug(string ...interface{})
		Info(string ...interface{})
		Print(string ...interface{})
		Warn(string ...interface{})
		Warning(string ...interface{})
		Error(string ...interface{})
		Fatal(string ...interface{})
		Panic(string ...interface{})
	}

	ResourceService interface {
		GetAddonScript() string
		GetCreateScript(path string) (string, error)
	}

	PackService interface {
		DescribeByReference(reference reference.Reference) (*rocketpack.RocketPack, error)
		FindByReference(ref reference.Reference) (*rocketpack.RocketPack, error)
		InstallByReference(reference reference.Reference, force bool) error
		UninstallByReference(reference reference.Reference) error
	}

	Driver struct {
		log              Logger          // the logger rocketblend will use for logging
		resource         ResourceService // the resource service rocketblend will use
		pack             PackService     // the pack service rocketblend will use
		debug            bool
		platform         runtime.Platform
		defaultBuild     string
		installationPath string
		addonsEnabled    bool
	}

	Options struct {
		Logger           Logger
		Debug            bool
		Platform         runtime.Platform
		InstallationPath string
		AddonsEnabled    bool
	}
)
