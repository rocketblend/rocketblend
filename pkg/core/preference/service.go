package preference

import "fmt"

type (
	Repository interface {
		Find() (*Settings, error)
		Create(i *Settings) error
	}

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

func (s *Service) Find() (*Settings, error) {
	conf, err := s.repo.Find()
	if err != nil {
		err = s.Create(&Settings{})
		if err != nil {
			return nil, fmt.Errorf("settings not found: %w", err)
		}
	}

	return conf, err
}

func (s *Service) Create(c *Settings) error {
	if err := s.repo.Create(c); err != nil {
		return fmt.Errorf("failed to insert settings: %w", err)
	}

	return nil
}
