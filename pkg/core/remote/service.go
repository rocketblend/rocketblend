package remote

import "fmt"

type (
	Repository interface {
		FindAll() ([]*Remote, error)
		Create(r *Remote) error
		Remove(name string) error
	}

	Config struct {
	}

	// Service is an remote Service
	Service struct {
		conf Config
		repo Repository
	}
)

// Create a new service
func NewService(conf Config, r Repository) *Service {
	srv := &Service{
		conf: conf,
		repo: r,
	}

	return srv
}

func LoadConfig() Config {
	return Config{}
}

// FindAll return all remotes
func (s *Service) FindAll() ([]*Remote, error) {
	remotes, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to find remotes: %w", err)
	}

	return remotes, nil
}

// Add a remote
func (s *Service) Add(r *Remote) error {
	if err := s.repo.Create(r); err != nil {
		return fmt.Errorf("failed to insert remote: %w", err)
	}

	return nil
}

// Remove a remote
func (s *Service) Remove(name string) error {
	if err := s.repo.Remove(name); err != nil {
		return fmt.Errorf("failed to remote remote: %w", err)
	}

	return nil
}
