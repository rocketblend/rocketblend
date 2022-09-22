package addon

import "fmt"

type (
	Repository interface {
		FindAll() ([]*Addon, error)
		FindByID(id string) (*Addon, error)
		Create(i *Addon) error
		Remove(id string) error
	}

	// Service is an addon Service
	Service struct {
		repo Repository
	}
)

func NewService(r Repository) *Service {
	srv := &Service{
		repo: r,
	}

	return srv
}

// FindAll return all addons
func (s *Service) FindAll() ([]*Addon, error) {
	addons, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to find addons: %w", err)
	}

	return addons, nil
}

// FindById return an addon by ID
func (s *Service) FindByID(id string) (*Addon, error) {
	addon, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find addon: %w", err)
	}

	return addon, err
}

func (s *Service) Create(i *Addon) error {
	if err := s.repo.Create(i); err != nil {
		return fmt.Errorf("failed to insert addon: %w", err)
	}

	return nil
}

func (s *Service) Remove(source string) error {
	if err := s.repo.Remove(source); err != nil {
		return fmt.Errorf("failed to remove addon: %w", err)
	}

	return nil
}
