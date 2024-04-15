package types

import "github.com/flowshot-io/x/pkg/logger"

type (
	Factory interface {
		GetLogger() (logger.Logger, error)
		GetRepository() (Repository, error)
		GetBlender() (Blender, error)
	}
)
