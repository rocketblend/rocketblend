package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// validateBlendFile checks if the path is non-empty and ends with .blend
func ValidateBlendFile(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if path == "" || !strings.HasSuffix(path, ".blend") {
		return false
	}

	return true
}
