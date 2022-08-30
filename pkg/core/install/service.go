package install

import "fmt"

type (
	Repository interface {
		FindAll() ([]*Install, error)
		FindByID(id string) (*Install, error)
		Create(i *Install) error
		Remove(id string) error
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
func (s *Service) FindAll() ([]*Install, error) {
	installs, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to find installs: %w", err)
	}

	return installs, nil
}

// FindById return an install by ID
func (s *Service) FindByID(id string) (*Install, error) {
	install, err := s.repo.FindByID(id)
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

func (s *Service) Remove(source string) error {
	if err := s.repo.Remove(source); err != nil {
		return fmt.Errorf("failed to remove install: %w", err)
	}

	return nil
}
