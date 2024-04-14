package blender

import (
	"fmt"
	"regexp"
	"strconv"
)

const renderOutputPattern = `Fra:(\d+) Mem:([0-9.]+[MK]?) \(Peak ([0-9.]+[MK]?)\) \| Time:([0-9:.]+) \| (.*)`

type (
	renderInfo struct {
		FrameNumber int
		Memory      string
		PeakMemory  string
		Time        string
		Operation   string
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

	return &renderInfo{
		FrameNumber: frameNumber,
		Memory:      match[2],
		PeakMemory:  match[3],
		Time:        match[4],
		Operation:   match[5],
	}, nil
}
