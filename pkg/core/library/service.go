package library

type (
	HTTPClient interface {
		FetchBuild(str string) (*Build, error)
		FetchPackage(str string) (*Package, error)
	}

	Repository interface {
		CreateBuild(build *Build) error
		CreatePackage(pack *Package) error
		FindBuildByRef(ref string) (*Build, error)
		FindPackageByRef(ref string) (*Package, error)
	}

	Service struct {
		client HTTPClient
		repo   Repository
	}
)

func NewService(client HTTPClient, repo Repository) *Service {
	srv := &Service{
		client: client,
		repo:   repo,
	}

	return srv
}

func (s *Service) FindBuildByRef(ref string) (*Build, error) {
	b, err := s.repo.FindBuildByRef(ref)
	if err != nil {
		b, err = s.fetchBuild(ref)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (s *Service) FindPackageByRef(ref string) (*Package, error) {
	p, err := s.repo.FindPackageByRef(ref)
	if err != nil {
		p, err = s.fetchPackage(ref)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (s *Service) fetchBuild(ref string) (*Build, error) {
	b, err := s.client.FetchBuild(ref)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateBuild(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) fetchPackage(ref string) (*Package, error) {
	p, err := s.client.FetchPackage(ref)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreatePackage(p); err != nil {
		return nil, err
	}

	return p, nil
}
