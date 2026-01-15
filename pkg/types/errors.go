package types

import "errors"

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFileExists   = errors.New("file already exists")

	ErrMissingBlenderBuild = errors.New("missing blender build")

	ErrUnsafeDomain = errors.New("download blocked: domain not in safe list")
)
