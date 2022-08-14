package build

type (
	Http interface {
		Fetch(remote string, platform string, tag string) ([]*Build, error)
	}

	Config struct {
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
		builds, err := s.http.Fetch(remote.URL, req.Platform, req.Tag)
		if err != nil {
			return nil, err
		}

		availableBuilds = append(availableBuilds, builds...)
	}

	return availableBuilds, nil
}

func (s *Service) Find(hash string) (*Build, error) {
	return nil, nil
}
