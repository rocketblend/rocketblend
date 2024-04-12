package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rocketblend/rocketblend/pkg/types"
)

// validateBlendFile checks if the path is non-empty and ends with .blend
func ValidateBlendFile(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if path == "" || !strings.HasSuffix(path, ".blend") {
		return false
	}

	return true
}

// ValidateOneBuild checks if there is only one build type dependency
func ValidateOneBuild(fl validator.FieldLevel) bool {
	deps, ok := fl.Field().Interface().([]*types.Installation)
	if !ok {
		return false
	}

	buildCount := 0
	for _, dep := range deps {
		if dep.Type == "build" {
			buildCount++
		}
		if buildCount > 1 {
			return false
		}
	}

	return buildCount == 1
}
