package rocketpack

import "github.com/flowshot-io/x/pkg/logger"

type LoggerWriter struct {
	Logger logger.Logger
}

func (lw LoggerWriter) Write(p []byte) (n int, err error) {
	lw.Logger.Info(string(p))
	return len(p), nil
}
