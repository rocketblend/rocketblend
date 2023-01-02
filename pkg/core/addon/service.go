package addon

import (
	"encoding/json"
	"net/url"

	"github.com/rocketblend/rocketblend/pkg/jot"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

const PackgeFile = "package.json"

type Service struct {
	driver *jot.Driver
}

func NewService(driver *jot.Driver) *Service {
	return &Service{
		driver: driver,
	}
}

func (srv *Service) FindByReference(ref reference.Reference) (*Package, error) {
	b, err := srv.driver.Read(ref, PackgeFile)
	if err != nil {
		return nil, err
	}

	p := &Package{}
	if err := json.Unmarshal(b, p); err != nil {
		return nil, err
	}

	return p, err
}

func (srv *Service) FetchByReference(ref reference.Reference) error {
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

func (srv *Service) PullByReference(ref reference.Reference) error {
	pack, err := srv.FindByReference(ref)
	if err != nil {
		return err
	}

	err = srv.driver.Write(ref, pack.Source.File, pack.Source.URL)
	if err != nil {
		return err
	}

	return nil
}
