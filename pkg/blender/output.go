package blender

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	renderOutputPattern = `Fra:(\d+) Mem:([0-9.]+[MK]?) \(Peak ([0-9.]+[MK]?)\) \| Time:([0-9:.]+) \| (.*)`

	generalOperation   = "blender"
	syncingOperation   = "syncing"
	renderingOperation = "rendering"
)

type (
	renderInfo struct {
		FrameNumber int
		Memory      string
		PeakMemory  string
		Time        string
		Operation   string
		Data        map[string]string
	}
)

func parseRenderOutput(line string) (*renderInfo, error) {
	re := regexp.MustCompile(renderOutputPattern)
	match := re.FindStringSubmatch(line)

	if len(match) != 6 {
		return nil, fmt.Errorf("could not parse line: %s", line)
	}

	frameNumber, err := strconv.Atoi(match[1])
	if err != nil {
		return nil, err
	}

	operationRaw := strings.ToLower(match[5])
	operationType, operationDetails := parseOperation(operationRaw)

	info := &renderInfo{
		FrameNumber: frameNumber,
		Memory:      strings.ToLower(match[2]),
		PeakMemory:  strings.ToLower(match[3]),
		Time:        strings.ToLower(match[4]),
		Operation:   operationType,
		Data:        make(map[string]string),
	}

	if operationType == renderingOperation {
		currentSample, totalSamples := parseSamples(operationDetails)
		info.Data["progress"] = strconv.Itoa(currentSample)
		info.Data["total"] = strconv.Itoa(totalSamples)
	}

	if operationType == syncingOperation {
		info.Data["object"] = operationDetails
	}

	return info, nil
}

func parseSamples(details string) (currentSample int, totalSamples int) {
	re := regexp.MustCompile(`(\d+) / (\d+) samples`)
	match := re.FindStringSubmatch(details)
	if len(match) == 3 {
		currentSample, _ = strconv.Atoi(match[1])
		totalSamples, _ = strconv.Atoi(match[2])
	}
	return
}

func parseOperation(operation string) (string, string) {
	operations := map[string]bool{
		syncingOperation:   true,
		renderingOperation: true,
	}

	for key := range operations {
		if strings.Contains(operation, key) {
			return key, strings.TrimPrefix(operation, key+" ")
		}
	}

	return generalOperation, operation
}
