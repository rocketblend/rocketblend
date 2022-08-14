package build

import (
	"fmt"

	"github.com/rocketblend/rocketblend/pkg/core/remote"
)

type (
	Http interface {
		Fetch(remote string, platform string, tag string) ([]*Build, error)
	}

	Config struct {
		Platform string
	}

	// Service is an install Service
	Service struct {
		conf Config
		http Http
	}
)

func NewService(conf Config, http Http) *Service {
	srv := &Service{
		conf: conf,
		http: http,
	}

	return srv
}

func NewConfig() Config {
	return Config{}
}

func (s *Service) FetchAll(req FetchRequest) ([]*Build, error) {
	var availableBuilds []*Build

	for _, remote := range req.Remotes {
		builds, err := s.http.Fetch(remote.URL, s.conf.Platform, req.Tag)
		if err != nil {
			return nil, err
		}

		availableBuilds = append(availableBuilds, builds...)
	}

	return availableBuilds, nil
}

func (s *Service) Find(remotes []*remote.Remote, hash string) (*Build, error) {
	// TODO: Update remotes to require lookup by hash endpoint.
	// Temporary workaround: Find by hash.

	builds, err := s.FetchAll(FetchRequest{Remotes: remotes})
	if err != nil {
		return nil, err
	}

	for _, build := range builds {
		if build.Hash == hash {
			return build, nil
		}
	}

	return nil, fmt.Errorf("build not found")
}
