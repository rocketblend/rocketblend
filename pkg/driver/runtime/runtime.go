package runtime

import "runtime"

func DetectPlatform() Platform {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "linux":
		return Linux
	case "darwin":
		if runtime.GOARCH == "amd64" {
			return DarwinAmd
		}
		return DarwinArm
	}
	return Undefined
}
