package types

import (
	"github.com/rocketblend/rocketblend/pkg/reference"
	"github.com/rocketblend/rocketblend/pkg/runtime"
)

const (
	DefaultBuild = "github.com/rocketblend/official-library/packages/v0/builds/blender/4.2.2"

	UntrustedDomainPolicyBlock UntrustedDomainPolicy = "block"
	UntrustedDomainPolicyWarn  UntrustedDomainPolicy = "warn"
	UntrustedDomainPolicyAllow UntrustedDomainPolicy = "allow"
)

var (
	DefaultAliases = map[string]string{
		"github.com/rocketblend/official-library/packages/v0/builds":         "builds",
		"github.com/rocketblend/official-library/packages/v0/addons":         "addons",
		"github.com/rocketblend/official-library/packages/v0/builds/blender": "blender",
	}

	DefaultTrustedDomains = []string{
		"download.blender.org",
	}
)

type (
	UntrustedDomainPolicy string

	Config struct {
		Platform              runtime.Platform      `mapstructure:"platform"`
		DefaultBuild          reference.Reference   `mapstructure:"defaultBuild"`
		InstallationsPath     string                `mapstructure:"installationsPath"`
		PackagesPath          string                `mapstructure:"packagesPath"`
		Aliases               map[string]string     `mapstructure:"aliases"`
		TrustedDomains        []string              `mapstructure:"trustedDomains"`
		UntrustedDomainPolicy UntrustedDomainPolicy `mapstructure:"untrustedDomainPolicy"`
	}

	Configurator interface {
		Get() (config *Config, err error)
		GetAllValues() map[string]interface{}
		GetValueByString(key string) string
		SetValueByString(key string, value string) error
		Save(config *Config) error
		Path() string
	}
)
