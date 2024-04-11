package types

import (
	"github.com/rocketblend/rocketblend/pkg/driver/reference"
)

type (
	Dependency struct {
		Reference reference.Reference `json:"reference" validate:"required"`
		Type      PackageType         `json:"type,omitempty" validate:"omitempty oneof=build addon"`
	}
)
