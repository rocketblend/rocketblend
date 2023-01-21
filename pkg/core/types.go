package core

import (
	"github.com/rocketblend/rocketblend/pkg/core/config"
	"github.com/rocketblend/rocketblend/pkg/core/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
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
		FindByReference(ref reference.Reference) (*rocketpack.RocketPack, error)
		FetchByReference(ref reference.Reference) error
		PullByReference(ref reference.Reference) error
	}

	Driver struct {
		conf     *config.Config  // the config rocketblend will use
		log      Logger          // the logger rocketblend will use for logging
		resource ResourceService // the resource service rocketblend will use
		pack     PackService     // the pack service rocketblend will use
	}

	Options struct {
		Config          *config.Config // the config rocketblend will use (configurable)
		Logger                         // the logger jot will use (configurable)
		ResourceService                // the resource service rocketblend will use (configurable)
		PackService                    // the pack service rocketblend will use (configurable)
	}
)
