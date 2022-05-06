package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

type LoggerObject interface {
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
}

func GetLogger() LoggerObject {
	if Logger == nil {
		Logger, _ = zap.NewProduction()
	}
	return Logger.Sugar()
}
