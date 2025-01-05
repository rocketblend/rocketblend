package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/types"
)

const (
	savedFilePattern = `(?i)^saved: '(.+)'$`
	quitPattern      = `^blender quit$`

	eventPatternEevee  = `Fra:(\d+) Mem:([0-9.]+[MK]?) \(Peak ([0-9.]+[MK]?)\) \| Time:([0-9:.]+) \| (.*)`
	eventPatternCycles = `Fra:(\d+) Mem:([0-9.]+[MK]?) \(Peak ([0-9.]+[MK]?)\) \| Time:([0-9:.]+) \| Mem:([0-9.]+[MK]?), Peak:([0-9.]+[MK]?) \| (.*)`
)

// ParseBlenderEvent parses Blender output and returns a BlenderEvent.
func ParseBlenderEvent(output string) (types.BlenderEvent, error) {
	if output == "" {
		return nil, fmt.Errorf("output is empty")
	}

	if event, err := parseQuitEvent(output); err == nil {
		return event, nil
	}

	if event, err := parseSavedFileEvent(output); err == nil {
		return event, nil
	}

	if event, err := parseEeveeBlenderEvent(output); err == nil {
		return event, nil
	}

	if event, err := parseCyclesBlenderEvent(output); err == nil {
		return event, nil
	}

	return nil, errors.New("could not parse output")
}

func parseQuitEvent(line string) (types.BlenderEvent, error) {
	if strings.ToLower(strings.TrimSpace(line)) == "blender quit" {
		return &types.QuitEvent{}, nil
	}

	return nil, fmt.Errorf("not a quit event")
}

func parseSavedFileEvent(line string) (types.BlenderEvent, error) {
	re := regexp.MustCompile(savedFilePattern)
	match := re.FindStringSubmatch(strings.TrimSpace(line))
	if len(match) != 2 {
		return nil, fmt.Errorf("could not parse saved file line: %s", line)
	}

	return &types.SavedFileEvent{
		Path: match[1],
	}, nil
}

func parseEeveeBlenderEvent(line string) (types.BlenderEvent, error) {
	return parseRenderEventWithPattern(line, eventPatternEevee)
}

func parseCyclesBlenderEvent(line string) (types.BlenderEvent, error) {
	return parseRenderEventWithPattern(line, eventPatternCycles)
}

func parseRenderEventWithPattern(line, pattern string) (types.BlenderEvent, error) {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(line)

	if len(match) < 6 {
		return nil, fmt.Errorf("could not parse line: %s", line)
	}

	frameNumber, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, err
	}

	operationRaw := strings.ToLower(match[len(match)-1])
	base := types.RenderBase{
		Frame:      frameNumber,
		Memory:     strings.ToLower(match[2]),
		PeakMemory: strings.ToLower(match[3]),
		Time:       strings.ToLower(match[4]),
	}

	return createRenderEventFromOperation(operationRaw, base)
}

func createRenderEventFromOperation(operationRaw string, base types.RenderBase) (types.BlenderEvent, error) {
	for op, handler := range operationRegistry {
		if strings.Contains(operationRaw, op) {
			return handler(operationRaw, base), nil
		}
	}

	return nil, nil
}
