package library

type (
	Http interface {
		FetchBuild(str string) (*Build, error)
		FetchPackage(str string) (*Package, error)
	}

	Repository interface {
		FindBuildByPath(path string) (*Build, error)
		FindPackageByPath(path string) (*Build, error)
	}

	Service struct {
		http Http
		repo Repository
	}
)

func NewService(http Http, repo Repository) *Service {
	srv := &Service{
		http: http,
		repo: repo,
	}

	return srv
}

func (s *Service) FindBuildByPath(path string) (*Build, error) {
	b, err := s.repo.FindBuildByPath(path)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) FindPackageByPath(path string) (*Build, error) {
	p, err := s.repo.FindPackageByPath(path)
	if err != nil {
		return nil, err
	}

	return p, nil
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
