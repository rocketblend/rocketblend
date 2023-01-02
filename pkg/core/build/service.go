package build

import (
	"encoding/json"
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/rocketblend/pkg/jot"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

const BuildFile = "build.json"

type (
	AddonService interface {
		FetchByReference(ref reference.Reference) error
		PullByReference(ref reference.Reference) error
	}

	Service struct {
		driver       *jot.Driver
		addonService AddonService
	}
)

func NewService(driver *jot.Driver, addonService AddonService) *Service {
	return &Service{
		driver:       driver,
		addonService: addonService,
	}
}

func (srv *Service) FindByReference(ref reference.Reference) (*Build, error) {
	b, err := srv.driver.Read(ref, BuildFile)
	if err != nil {
		return nil, err
	}

	build := &Build{}
	if err := json.Unmarshal(b, build); err != nil {
		return nil, err
	}

	return build, err
}

func (srv *Service) FetchByReference(ref reference.Reference) error {
	err := srv.driver.Write(ref, BuildFile, ref.Url())
	if err != nil {
		return err
	}

	build, err := srv.FindByReference(ref)
	if err != nil {
		return err
	}

	for _, pack := range build.Packages {
		err = srv.addonService.FetchByReference(reference.Reference(pack))
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *Service) PullByReference(ref reference.Reference, platform runtime.Platform) error {
	build, err := srv.FindByReference(ref)
	if err != nil {
		return err
	}

	source := build.GetSourceForPlatform(platform)
	if source == nil {
		return fmt.Errorf("no source found for platform %s", platform)
	}

	err = srv.driver.WriteAndExtract(ref, jot.GetFilenameFromURL(source.URL), source.URL)
	if err != nil {
		return err
	}

	for _, pack := range build.Packages {
		err = srv.addonService.PullByReference(reference.Reference(pack))
		if err != nil {
			return err
		}
	}

	return nil
}
