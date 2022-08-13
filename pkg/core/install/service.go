package install

import "fmt"

type (
	Repository interface {
		FindAll(req FindRequest) ([]*Install, error)
		FindByHash(hash string) (*Install, error)
		Create(i *Install) error
		Remove(hash string) error
	}

	Installer interface {
		Download(url string) error
	}

	Config struct {
	}

	// Service is an install Service
	Service struct {
		conf Config
		repo Repository
	}
)

func NewService(conf Config, r Repository) *Service {
	srv := &Service{
		conf: conf,
		repo: r,
	}

	return srv
}

func LoadConfig() Config {
	var conf Config
	return conf
}

// FindAll return all installs
func (s *Service) FindAll(req FindRequest) ([]*Install, error) {
	installs, err := s.repo.FindAll(req)
	if err != nil {
		return nil, fmt.Errorf("failed to find installs: %w", err)
	}

	return installs, nil
}

// FindByHash return an install by hash
func (s *Service) FindByHash(hash string) (*Install, error) {
	install, err := s.repo.FindByHash(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to find install: %w", err)
	}

	return install, err
}

func (s *Service) Create(i *Install) error {
	if err := s.repo.Create(i); err != nil {
		return fmt.Errorf("failed to insert install: %w", err)
	}

	return nil
}

func (s *Service) Remove(hash string) error {
	if err := s.repo.Remove(hash); err != nil {
		return fmt.Errorf("failed to remove install: %w", err)
	}

	return nil
}
