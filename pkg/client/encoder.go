package client

import "github.com/rocketblend/rocketblend/pkg/core/encoder"

func NewEncoderService() *encoder.Service {
	srv := encoder.NewService()

	return srv
}
