package blenderparser

import (
	"fmt"
	"regexp"
	"strconv"
)

type RenderInfo struct {
	FrameNumber int
	Memory      string
	PeakMemory  string
	Time        string
	Operation   string
}

func ParseRenderOutput(line string) (RenderInfo, error) {
	re := regexp.MustCompile(`Fra:(\d+) Mem:([0-9.]+[MK]?) \(Peak ([0-9.]+[MK]?)\) \| Time:([0-9:.]+) \| (.*)`)
	match := re.FindStringSubmatch(line)

	if len(match) != 6 {
		return RenderInfo{}, fmt.Errorf("could not parse line: %s", line)
	}

	frameNumber, err := strconv.Atoi(match[1])
	if err != nil {
		return RenderInfo{}, err
	}

	info := RenderInfo{
		FrameNumber: frameNumber,
		Memory:      match[2],
		PeakMemory:  match[3],
		Time:        match[4],
		Operation:   match[5],
	}

	return info, nil
}
