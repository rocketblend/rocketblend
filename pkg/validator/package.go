package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/types"
)

func PackageDependenciesValidator(sl validator.StructLevel) {
	rocketPack, ok := sl.Current().Interface().(types.Package)
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

func ValidateUniquePlatforms(sl validator.StructLevel) {
	rp := sl.Current().Interface().(types.Package)

	platformCount := make(map[types.Platform]int)
	for _, source := range rp.Sources {
		if source != nil {
			platformCount[source.Platform]++
		}
	}

	// Check for duplicates including "" and "any"
	for platform, count := range platformCount {
		if count > 1 {
			var errorTag string
			if platform == "" {
				errorTag = "UniqueEmptyPlatform"
			} else if platform == "any" {
				errorTag = "UniqueAnyPlatform"
			} else {
				errorTag = "UniquePlatform"
			}

			sl.ReportError(rp.Sources, "Sources", "Sources", errorTag, fmt.Sprintf("multiple entries for platform %s", platform))
		}
	}
}
