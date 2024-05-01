package helpers

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ExtractName(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func PadWithZero(num int, padLength int) string {
	numStr := fmt.Sprintf("%d", num)
	if len(numStr) >= padLength {
		return numStr
	}

	return strings.Repeat("0", padLength-len(numStr)) + numStr // Pad and return
}
