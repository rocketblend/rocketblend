package client

import (
	"github.com/rocketblend/rocketblend/pkg/core/preference"
	"github.com/rocketblend/scribble"
)

func NewPreferenceService(driver *scribble.Driver) *preference.Service {
	repo := preference.NewScribbleRepository(driver)
	srv := preference.NewService(repo)

	return srv
}
