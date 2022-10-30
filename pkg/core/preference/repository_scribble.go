package preference

import (
	"fmt"

	"github.com/rocketblend/scribble"
)

type ScribbleRepository struct {
	driver     *scribble.Driver
	collection string
	name       string
}

func NewScribbleRepository(driver *scribble.Driver) *ScribbleRepository {
	return &ScribbleRepository{
		driver:     driver,
		collection: "preferences",
		name:       "settings",
	}
}

func (r *ScribbleRepository) Find() (*Settings, error) {
	settings := Settings{}
	if err := r.driver.Read(r.collection, r.name, &settings); err != nil {
		return nil, fmt.Errorf("preferences not found: %s", r.name)
	}

	return &settings, nil
}

func (r *ScribbleRepository) Create(i *Settings) error {
	if err := r.driver.Write(r.collection, r.name, i); err != nil {
		return err
	}

	return nil
}
