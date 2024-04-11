package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/driver/types"
)

func DependenciesValidation(sl validator.StructLevel) {
	deps := sl.Current().Interface().(types.Dependencies)

	// Validate only one build type dependency in Direct dependencies.
	buildTypeCount := 0
	for _, dep := range deps.Direct {
		if dep.Type == types.PackageBuild {
			buildTypeCount++
		}
	}
	if buildTypeCount > 1 {
		sl.ReportError(deps.Direct, "Direct", "Direct", "OnlyOneBuildType", "")
	}

	// Validate Indirect dependencies can exist only if there's at least one Direct dependency.
	if len(deps.Indirect) > 0 && len(deps.Direct) == 0 {
		sl.ReportError(deps.Indirect, "Indirect", "Indirect", "IndirectWithNoDirect", "")
	}
}
