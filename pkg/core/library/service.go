package library

type (
	Http interface {
		FetchBuild(url string) (*Build, error)
		FetchPackage(url string) (*Package, error)
	}

	Service struct {
		http Http
	}
)

func NewService(http Http) *Service {
	srv := &Service{
		http: http,
	}

	return srv
}

func (s *Service) FetchBuild(str string) (*Build, error) {
	b, err := s.http.FetchBuild(str)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) FetchPackage(str string) (*Package, error) {
	p, err := s.http.FetchPackage(str)
	if err != nil {
		return nil, err
	}

	return p, nil
}
