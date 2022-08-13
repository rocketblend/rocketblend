package install

import (
	"encoding/json"

	"github.com/rocketblend/rocketblend/pkg/scribble"
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

func (r *ScribbleRepository) FindAll(req FindRequest) ([]*Install, error) {
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

func (r *ScribbleRepository) FindByHash(hash string) (*Install, error) {
	var install *Install
	if err := r.driver.Read(r.collection, hash, install); err != nil {
		return nil, err
	}

	return install, nil
}

func (r *ScribbleRepository) Create(i *Install) error {
	if err := r.driver.Write(r.collection, i.Hash, i); err != nil {
		return err
	}

	return nil
}

func (r *ScribbleRepository) Remove(hash string) error {
	if err := r.driver.Delete(r.collection, hash); err != nil {
		return err
	}

	return nil
}
