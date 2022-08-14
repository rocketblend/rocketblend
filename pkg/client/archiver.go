package client

import "github.com/rocketblend/rocketblend/pkg/core/archiver"

func NewArchiverService(delete bool) *archiver.Service {
	conf := archiver.NewConfig(delete)
	srv := archiver.NewService(conf)

	return srv
}
