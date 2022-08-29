package install

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
		collection: "installs",
	}
}

func (r *ScribbleRepository) FindAll() ([]*Install, error) {
	var installs []*Install

	records, err := r.driver.ReadAll(r.collection)
	if err != nil {
		return nil, err
	}

	for _, i := range records {
		installFound := Install{}
		if err := json.Unmarshal([]byte(i), &installFound); err != nil {
			return nil, err
		}
		installs = append(installs, &installFound)
	}

	return installs, nil
}

func (r *ScribbleRepository) FindBySource(source string) (*Install, error) {
	install := Install{}
	if err := r.driver.Read(r.collection, source, &install); err != nil {
		return nil, fmt.Errorf("install not found: %s", source)
	}

	return &install, nil
}

func (r *ScribbleRepository) Create(i *Install) error {
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
