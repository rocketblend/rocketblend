package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/types"
)

type (
	Validator struct {
		validator *validator.Validate
	}
)

func New() *Validator {
	validate := validator.New(
		validator.WithRequiredStructEnabled(),
	)

	validate.RegisterValidation("blendfile", ValidateBlendFile)
	validate.RegisterValidation("onebuild", ValidateOneBuild)

	validate.RegisterStructValidation(DependenciesValidation, types.Dependencies{})
	validate.RegisterStructValidation(PackageDependenciesValidator, types.Package{})
	validate.RegisterStructValidation(ValidateUniquePlatforms, types.Package{})

	return &Validator{
		validator: validate,
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
