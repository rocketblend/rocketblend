package runtime

type Platform int64

const (
	Undefined Platform = iota
	Windows
	Linux
	DarwinAmd
	DarwinArm
)

func (p Platform) String() string {
	switch p {
	case Windows:
		return "windows"
	case Linux:
		return "linux"
	case DarwinAmd:
		return "macos/intel"
	case DarwinArm:
		return "macos/apple"
	}
	return "unknown"
}
