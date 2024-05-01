package blender

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	RevisionTempalteVariable = "{{.Revision}}"
	NameTemplateVariable     = "{{.Name}}"
)

type (
	TemplatedOutputData struct {
		Name     string `json:"name"`
		Revision string `json:"revision"`
	}
)

var ErrInvalidPath = errors.New("path should not start with //")

func FindMaxRevision(templatedPath string) (int, error) {
	if strings.HasPrefix(templatedPath, "//") {
		return -1, ErrInvalidPath
	}

	dirPath := filepath.Dir(templatedPath)

	revisionDirPattern := regexp.MustCompile(RevisionTempalteVariable)
	if !revisionDirPattern.MatchString(dirPath) {
		return -1, fmt.Errorf("no {{.Revision}} placeholder found in the path")
	}

	searchPath := revisionDirPattern.ReplaceAllString(dirPath, "*")

	directories, err := filepath.Glob(searchPath)
	if err != nil {
		return -1, err
	}

	maxRevision := -1
	revisionNumberPattern := regexp.MustCompile(`\d+$`)
	for _, dir := range directories {
		base := filepath.Base(dir)
		if matches := revisionNumberPattern.FindStringSubmatch(base); matches != nil {
			revNum, _ := strconv.Atoi(matches[0])
			if revNum > maxRevision {
				maxRevision = revNum
			}
		}
	}

	return maxRevision, nil
}

func FindMaxFrameNumber(templatedPath string) (int, error) {
	baseName := filepath.Base(templatedPath)
	directory := filepath.Dir(templatedPath)

	numFormatRegex := regexp.MustCompile(`#+`)
	numFormat := numFormatRegex.FindString(baseName)
	if numFormat == "" {
		return -1, fmt.Errorf("no number placeholder found in filename pattern")
	}

	regexPattern := regexp.QuoteMeta(strings.Replace(baseName, numFormat, "REPLACEME", 1))
	regexPattern = strings.Replace(regexPattern, "REPLACEME", `(\d{`+strconv.Itoa(len(numFormat))+`})`, 1)
	re := regexp.MustCompile(regexPattern)

	files, err := os.ReadDir(directory)
	if err != nil {
		return -1, fmt.Errorf("unable to read directory %s: %w", directory, err)
	}

	maxFrame := -1

	for _, file := range files {
		if matches := re.FindStringSubmatch(file.Name()); matches != nil {
			frameNumber, err := strconv.Atoi(matches[1])
			if err != nil {
				continue // skip files with non-integer frame numbers
			}
			if frameNumber > maxFrame {
				maxFrame = frameNumber
			}
		}
	}

	if maxFrame == -1 {
		return -1, fmt.Errorf("no frame files found in %s", directory)
	}

	return maxFrame, nil
}
