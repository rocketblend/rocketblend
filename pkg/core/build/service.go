package build

import (
	"encoding/json"
	"fmt"
	"net/url"

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
		platform     runtime.Platform
	}
)

func NewService(driver *jot.Driver, platform runtime.Platform, addonService AddonService) *Service {
	return &Service{
		driver:       driver,
		platform:     platform,
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
	downloadUrl, err := url.JoinPath(ref.Url(), BuildFile)
	if err != nil {
		return err
	}

	err = srv.driver.Write(ref, BuildFile, downloadUrl)
	if err != nil {
		return err
	}

	build, err := srv.FindByReference(ref)
	if err != nil {
		return err
	}

	for _, pack := range build.Addons {
		err = srv.addonService.FetchByReference(reference.Reference(pack))
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *Service) PullByReference(ref reference.Reference) error {
	build, err := srv.FindByReference(ref)
	if err != nil {
		return err
	}

	source := build.GetSourceForPlatform(srv.platform)
	if source == nil {
		return fmt.Errorf("no source found for platform %s", (srv.platform))
	}

	err = srv.driver.WriteAndExtract(ref, jot.GetFilenameFromURL(source.URL), source.URL)
	if err != nil {
		return err
	}

	for _, pack := range build.Addons {
		err = srv.addonService.PullByReference(reference.Reference(pack))
		if err != nil {
			return err
		}
	}

	return nil
}
