package types

import "errors"

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFileExists   = errors.New("file already exists")

	ErrMissingBlenderBuild = errors.New("missing blender build")

	ErrUnsafeDomain    = errors.New("download blocked from unsafe domain")
	ErrUntrustedDomain = errors.New("download from untrusted domain requires user confirmation")
)
