package remote

import (
	"encoding/json"

	"github.com/rocketblend/scribble"
)

type ScribbleRepository struct {
	driver     *scribble.Driver
	collection string
}

func NewScribbleRepository(driver *scribble.Driver) *ScribbleRepository {
	return &ScribbleRepository{
		driver:     driver,
		collection: "remotes",
	}
}

func (r *ScribbleRepository) FindAll() ([]*Remote, error) {
	var remotes []*Remote

	records, err := r.driver.ReadAll(r.collection)
	if err != nil {
		return nil, err
	}

	for _, i := range records {
		remoteFound := Remote{}
		if err := json.Unmarshal([]byte(i), &remoteFound); err != nil {
			return nil, err
		}
		remotes = append(remotes, &remoteFound)
	}

	return remotes, nil
}

func (r *ScribbleRepository) Create(remote *Remote) error {
	if err := r.driver.Write(r.collection, remote.Name, remote); err != nil {
		return err
	}

	return nil
}

func (r *ScribbleRepository) Remove(name string) error {
	if err := r.driver.Delete(r.collection, name); err != nil {
		return err
	}

	return nil
}
