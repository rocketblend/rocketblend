package addon

import (
	"encoding/json"
	"fmt"

	"github.com/rocketblend/scribble"
)

type ScribbleRepository struct {
	driver     *scribble.Driver
	collection string
}

func NewScribbleRepository(driver *scribble.Driver) *ScribbleRepository {
	return &ScribbleRepository{
		driver:     driver,
		collection: "addons",
	}
}

func (r *ScribbleRepository) FindAll() ([]*Addon, error) {
	var addons []*Addon

	records, err := r.driver.ReadAll(r.collection)
	if err != nil {
		return nil, err
	}

	for _, i := range records {
		addonFound := Addon{}
		if err := json.Unmarshal([]byte(i), &addonFound); err != nil {
			return nil, err
		}
		addons = append(addons, &addonFound)
	}

	return addons, nil
}

func (r *ScribbleRepository) FindByID(id string) (*Addon, error) {
	addon := Addon{}
	if err := r.driver.Read(r.collection, id, &addon); err != nil {
		return nil, fmt.Errorf("addon not found: %s", id)
	}

	return &addon, nil
}

func (r *ScribbleRepository) Create(i *Addon) error {
	if err := r.driver.Write(r.collection, i.Id, i); err != nil {
		return err
	}

	return nil
}

func (r *ScribbleRepository) Remove(source string) error {
	if err := r.driver.Delete(r.collection, source); err != nil {
		return err
	}

	return nil
}
