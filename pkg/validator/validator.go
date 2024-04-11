package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Validate struct {
		*validator.Validate
	}
)

func New() *Validate {
	validate := validator.New(
		validator.WithRequiredStructEnabled(),
	)

	validate.RegisterStructValidation(DependenciesValidation, types.Dependencies{})
	validate.RegisterStructValidation(RocketPackDependenciesValidator, types.RocketPack{})

	return &Validate{validate}
}
