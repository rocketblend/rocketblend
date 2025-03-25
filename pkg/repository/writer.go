package repository

import "github.com/rocketblend/rocketblend/pkg/logger"

type LoggerWriter struct {
	Logger logger.Logger
}

func (lw LoggerWriter) Write(p []byte) (n int, err error) {
	lw.Logger.Debug("Git", map[string]interface{}{"message": string(p)})
	return len(p), nil
}
