package rocketblend2

import (
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendconfig"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/blendfile"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/installation"
	"github.com/rocketblend/rocketblend/pkg/rocketblend2/rocketpack"
)

type (
	RocketPackService interface {
		GetPackages(references ...reference.Reference) ([]*rocketpack.RocketPack, error)
		RemovePackages(references ...reference.Reference) error
	}

	InstallationService interface {
		GetInstallations(references ...reference.Reference) ([]*installation.Installation, error)
		RemoveInstallations(references ...reference.Reference) error
	}

	BlendConfigService interface {
		AddDependencies(config *blendconfig.BlendConfig, references ...reference.Reference) error
		RemoveDependencies(config *blendconfig.BlendConfig, references ...reference.Reference) error
		ResolveBlendFile(config *blendconfig.BlendConfig) (*blendfile.BlendFile, error)
	}

	BlendFileService interface {
		Render(blendFile *blendfile.BlendFile) error
		Run(blendFile *blendfile.BlendFile) error
		Create(blendFile *blendfile.BlendFile) error
	}

	RocketBlend interface {
		Render()
		Run()
		Create()

		AddDependencies(global bool, references ...reference.Reference) error
		RemoveDependencies(global bool, references ...reference.Reference) error
		DescribeDependencies(references ...reference.Reference) error
	}

	rocketBlend struct {
		rocketPackService   RocketPackService
		installationService InstallationService
		blendConfigService  BlendConfigService
		blendFileService    BlendFileService

		blendConfig blendconfig.BlendConfig
	}
)
