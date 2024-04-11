package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func RocketPackDependenciesValidator(sl validator.StructLevel) {
	rocketPack, ok := sl.Current().Interface().(types.RocketPack)
	if !ok {
		return
	}

	for _, dep := range rocketPack.Dependencies {
		if dep.Type == types.PackageBuild {
			sl.ReportError(dep.Type, "Dependencies", "Dependencies", "NoBuildDependenciesAllowed", "")
			break
		}
	}
}
