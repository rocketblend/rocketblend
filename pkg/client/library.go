package client

import "github.com/rocketblend/rocketblend/pkg/core/library"

func NewLibraryService() *library.Service {
	conf := library.NewClientConfig()
	http := library.NewClient(conf)
	repo := library.NewRepository()
	srv := library.NewService(http, repo)

	return srv
}
