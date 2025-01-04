package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/rocketblend/rocketblend/pkg/types"
)

const (
	eventPatternEevee  = `Fra:(\d+) Mem:([0-9.]+[MK]?) \(Peak ([0-9.]+[MK]?)\) \| Time:([0-9:.]+) \| (.*)`
	eventPatternCycles = `Fra:(\d+) Mem:([0-9.]+[MK]?) \(Peak ([0-9.]+[MK]?)\) \| Time:([0-9:.]+) \| Mem:([0-9.]+[MK]?), Peak:([0-9.]+[MK]?) \| (.*)`
)

// ParseBlenderEvent parses Blender output and returns a BlenderEvent.
func ParseBlenderEvent(output string) (types.BlenderEvent, error) {
	if output == "" {
		return nil, fmt.Errorf("output is empty")
	}

	if event, err := parseEeveeBlenderEvent(output); err == nil {
		return event, nil
	}

	if event, err := parseCyclesBlenderEvent(output); err == nil {
		return event, nil
	}

	return &types.GenericEvent{Message: output}, nil
}

// Parsing logic for Eevee Blender events
func parseEeveeBlenderEvent(line string) (types.BlenderEvent, error) {
	return parseRenderEventWithPattern(line, eventPatternEevee)
}

// Parsing logic for Cycles Blender events
func parseCyclesBlenderEvent(line string) (types.BlenderEvent, error) {
	return parseRenderEventWithPattern(line, eventPatternCycles)
}

// Common parsing logic shared by Eevee and Cycles
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
	event := &types.RenderEvent{
		Frame:      frameNumber,
		Memory:     strings.ToLower(match[2]),
		PeakMemory: strings.ToLower(match[3]),
		Time:       strings.ToLower(match[4]),
		Data:       make(map[string]string),
	}

	handleOperation(operationRaw, event)
	return event, nil
}
