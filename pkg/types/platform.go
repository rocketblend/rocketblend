package types

type Platform string

const (
	PlatformAny         Platform = "any"
	PlatformWindows     Platform = "windows"
	PlatformLinux       Platform = "linux"
	PlatformDarwinAMD   Platform = "macos/intel"
	PlatformDarwinARM   Platform = "macos/apple"
	PlatformUnsupported Platform = "unsupported"
)
