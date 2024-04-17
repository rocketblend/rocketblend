package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"unicode"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/pkg/errors"
)

func ValidateFilePath(filePath string, requiredFileName string) error {
	if filePath == "" {
		return errors.New("file path is empty")
	}

	if filepath.Base(filePath) != requiredFileName && requiredFileName != "" {
		return fmt.Errorf("invalid file name (must be '%s'): %s", requiredFileName, filepath.Base(filePath))
	}

	return nil
}

func FileExists(filePath string) error {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return errors.New("file does not exist")
	}

	if info.IsDir() {
		return errors.New("file is a directory")
	}

	return nil
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func LogAndReturnError(logger logger.Logger, msg string, err error, fields ...map[string]interface{}) error {
	var fieldMap map[string]interface{}
	if len(fields) > 0 {
		fieldMap = fields[0]
	} else {
		fieldMap = make(map[string]interface{})
	}
	fieldMap["error"] = err.Error()

	logger.Error(Capitalize(msg), fieldMap)
	return errors.Wrap(err, msg)
}
