package rocketpack

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"sigs.k8s.io/yaml"

	"github.com/rocketblend/rocketblend/pkg/core/runtime"
	"github.com/rocketblend/rocketblend/pkg/jot"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

const PackgeFile = "rocketpack.yaml"

type Service struct {
	driver   *jot.Driver
	platform runtime.Platform
}

func NewService(driver *jot.Driver, platform runtime.Platform) *Service {
	return &Service{
		driver:   driver,
		platform: platform,
	}
}

func (srv *Service) DescribeByReference(reference reference.Reference) (*RocketPack, error) {
	url, err := url.JoinPath(reference.Url(), PackgeFile)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pack := &RocketPack{}
	if err := yaml.Unmarshal(bodyBytes, pack); err != nil {
		return nil, err
	}

	err = validate(pack)
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (srv *Service) InstallByReference(ref reference.Reference, force bool) error {
	// Check if already installed.
	pack, _ := srv.FindByReference(ref)

	if pack == nil || force {
		err := srv.fetchByReference(ref)
		if err != nil {
			return err
		}

		err = srv.pullByReference(ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *Service) UninstallByReference(ref reference.Reference) error {
	_, err := srv.FindByReference(ref)
	if err != nil {
		return err
	}

	if err := srv.driver.DeleteAll(ref); err != nil {
		return err
	}

	return nil
}

func (srv *Service) FindByReference(ref reference.Reference) (*RocketPack, error) {
	b, err := srv.driver.Read(ref, PackgeFile)
	if err != nil {
		return nil, err
	}

	p := &RocketPack{}
	if err := yaml.Unmarshal(b, p); err != nil {
		return nil, err
	}

	err = validate(p)
	if err != nil {
		return nil, err
	}

	return p, err
}

func (srv *Service) fetchByReference(ref reference.Reference) error {
	// Validates reference is a valid pack.
	_, err := srv.DescribeByReference(ref)
	if err != nil {
		return err
	}

	downloadUrl, err := url.JoinPath(ref.Url(), PackgeFile)
	if err != nil {
		return err
	}

	err = srv.driver.Write(ref, PackgeFile, downloadUrl)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Service) pullByReference(ref reference.Reference) error {
	pack, err := srv.FindByReference(ref)
	if err != nil {
		return err
	}

	if pack.Addon != nil {
		return srv.writeAddon(ref, pack.Addon)
	}

	if pack.Build != nil {
		return srv.writeBuild(ref, pack.Build)
	}

	return fmt.Errorf("no build or addon found in rocketpack %s", ref)
}

func (srv *Service) writeAddon(ref reference.Reference, addon *Addon) error {
	err := srv.driver.Write(ref, addon.Source.File, addon.Source.URL)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Service) writeBuild(ref reference.Reference, build *Build) error {
	source := build.GetSourceForPlatform(srv.platform)
	if source == nil {
		return fmt.Errorf("no source found for platform %s", (srv.platform))
	}

	err := srv.driver.WriteAndExtract(ref, jot.GetFilenameFromURL(source.URL), source.URL)
	if err != nil {
		return err
	}

	for _, pack := range build.Addons {
		err = srv.pullByReference(reference.Reference(pack))
		if err != nil {
			return err
		}
	}

	return nil
}
