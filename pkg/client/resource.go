package client

import "github.com/rocketblend/rocketblend/pkg/core/resource"

func NewResourceService(dir string) *resource.Service {
	srv := resource.NewService(dir)
	return srv
}
