package parser

import (
	"regexp"
	"strconv"
)

func parseSamples(details string) (currentSample int, totalSamples int) {
	re := regexp.MustCompile(`(\d+)\s*/\s*(\d+)\s*(samples|sample)?`)
	match := re.FindStringSubmatch(details)
	if len(match) >= 3 {
		currentSample, _ = strconv.Atoi(match[1])
		totalSamples, _ = strconv.Atoi(match[2])
	}
	return
}
